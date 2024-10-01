package cplugin

import (
	"crypto/sha512"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"log"
	"log/slog"
	"os"
	"os/exec"

	"github.com/fsnotify/fsnotify"
)

func NewDiscovery[T any](conf *config) *discovery[T] {
	return &discovery[T]{
		conf:     conf,
		updateC:  make(chan string, 10),
		registry: NewRegistry[T](),
	}
}
func NewDiscoveryWithRegistry[T any](registry *registry[T], conf *config) *discovery[T] {
	if conf == nil {
		conf = NewDefaultConfig()
	}
	d := &discovery[T]{
		conf:     conf,
		updateC:  make(chan string, 10),
		registry: registry,
	}
	for _, preload := range conf.preloads {
		if !preload.lazy {
			d.updateC <- preload.filePath
		}
	}
	d.Listen()

	return d
}

type discovery[T any] struct {
	conf     *config
	updateC  chan string
	registry *registry[T]
}

func (d *discovery[T]) Destroy() error {
	return nil
}
func (d *discovery[T]) Listen() error {
	var stopC = make(chan struct{})
	go d.watch(stopC)
	go d.handleNewEvent(stopC)
	return nil
}
func (d *discovery[T]) watch(stopC chan struct{}) error {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	defer watcher.Close()
	err = watcher.Add(d.conf.discoverDirectory)
	if err != nil {
		return err
	}

	// Start listening for events
	for {
		select {
		case <-stopC:
			return nil
		case event, ok := <-watcher.Events:
			if !ok {
				return errors.New("watcher event error")
			}

			// Check if the event is a file creation
			if event.Op&fsnotify.Create == fsnotify.Create {

				// Add your file processing logic here
			}
			if event.Op&fsnotify.Write == fsnotify.Write {
				defaultLogger.Info("update event", slog.String("filename", event.Name))
				d.updateC <- event.Name
			}
		case err, ok := <-watcher.Errors:
			if !ok {
				return errors.New("watcher error")
			}
			defaultLogger.Error("watcher error", slog.String("err", err.Error()))
		}
	}
}
func (d *discovery[T]) handleNewEvent(stopC chan struct{}) error {
	for {
		select {
		case <-stopC:
			return nil
		case file := <-d.updateC:
			func() {
				mutex.Lock()
				defer mutex.Unlock()
				soPath, err := d.buildPlugin(file)
				if err != nil {
					defaultLogger.Warn("invalid plugin", slog.String("err", err.Error()))
					if err := os.Remove(file); err != nil {
						defaultLogger.Warn("remove invalid plugin fail", slog.String("err", err.Error()))
					}
				} else {
					pl, ok := getPluginFromCache(soPath)
					if ok {
						c, ok := pl.(T)
						if ok {
							d.registry.Add(file, c)
							return
						}
					}
					app, err := d.registry.AddByPath(file, soPath)
					if err != nil {
						defaultLogger.Error("add plugin to registry fail", slog.String("err", err.Error()), slog.String("name", file))
						return
					}

					addToCache(soPath, app)
					defaultLogger.Info("add plugin success", slog.String("name", file))
				}
			}()

		}
	}
}
func (d *discovery[T]) buildPlugin(name string) (string, error) {
	soPath := d.conf.genSoPath(name)

	goPath := d.conf.discoverDirectory + name
	cmdStr := fmt.Sprintf("go build -buildmode=plugin -gcflags=all='-N -l' -o %s %s", soPath, goPath)
	defaultLogger.Info("executing cmd", slog.String("cmd", cmdStr))
	cmd := exec.Command("bash", "-c", cmdStr)
	output, err := cmd.Output()
	if err != nil {
		defaultLogger.Error("execute cmd error: ", slog.String("cmd", cmdStr), slog.String("output", string(output)), slog.String("err", err.Error()))
		return "", err
	}

	return soPath, nil
}

func getPluginFromCache(path string) (any, bool) {
	id := getIdByPath(path)
	if len(id) < 1 {
		return "", false
	}
	v, ok := pluginCache[id]
	return v, ok
}
func getIdByPath(filePath string) string {
	id, err := hashFileWithSize(filePath)
	if err != nil {
		defaultLogger.Error("hashFileWithSize error: ", slog.String("path", filePath), slog.String("err", err.Error()))
		return ""
	}
	return id
}
func addToCache(filePath string, c any) {
	id, err := hashFileWithSize(filePath)
	if err != nil {
		defaultLogger.Error("hashFileWithSize error: ", slog.String("path", filePath), slog.String("err", err.Error()))
		return
	}
	pluginCache[id] = c
}
func hashFileWithSize(filePath string) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	// Get the file size
	fileInfo, err := file.Stat()
	if err != nil {
		return "", err
	}
	fileSize := fileInfo.Size()

	// Create a new SHA-512 hash
	hash := sha512.New()

	// Copy the file's content into the hash function
	if _, err := io.Copy(hash, file); err != nil {
		return "", err
	}

	// Get the final hash sum and encode it as a hexadecimal string
	hashInBytes := hash.Sum(nil)
	hashInString := hex.EncodeToString(hashInBytes)

	// Append the file size to the hash string
	result := fmt.Sprintf("%s_%d", hashInString, fileSize)

	return result, nil
}
