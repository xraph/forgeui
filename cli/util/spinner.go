package util

import (
	"fmt"
	"io"
	"os"
	"sync"
	"time"
)

// Spinner provides a simple loading spinner
type Spinner struct {
	message string
	frames  []string
	delay   time.Duration
	writer  io.Writer
	stop    chan bool
	wg      sync.WaitGroup
	active  bool
	mu      sync.Mutex
}

var defaultFrames = []string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"}

// NewSpinner creates a new spinner with a message
func NewSpinner(message string) *Spinner {
	return &Spinner{
		message: message,
		frames:  defaultFrames,
		delay:   100 * time.Millisecond,
		writer:  os.Stdout,
		stop:    make(chan bool),
	}
}

// Start starts the spinner animation
func (s *Spinner) Start() {
	s.mu.Lock()
	if s.active {
		s.mu.Unlock()
		return
	}
	s.active = true
	s.mu.Unlock()
	
	s.wg.Add(1)
	go s.animate()
}

// Stop stops the spinner animation
func (s *Spinner) Stop() {
	s.mu.Lock()
	if !s.active {
		s.mu.Unlock()
		return
	}
	s.active = false
	s.mu.Unlock()
	
	s.stop <- true
	s.wg.Wait()
	
	// Clear the line
	fmt.Fprintf(s.writer, "\r\033[K")
}

// Success stops the spinner and shows a success message
func (s *Spinner) Success(msg string) {
	s.Stop()
	fmt.Fprintf(s.writer, "%s✓%s %s\n", ColorGreen, ColorReset, msg)
}

// Error stops the spinner and shows an error message
func (s *Spinner) Error(msg string) {
	s.Stop()
	fmt.Fprintf(s.writer, "%s✗%s %s\n", ColorRed, ColorReset, msg)
}

// UpdateMessage updates the spinner message
func (s *Spinner) UpdateMessage(msg string) {
	s.mu.Lock()
	s.message = msg
	s.mu.Unlock()
}

// animate runs the spinner animation loop
func (s *Spinner) animate() {
	defer s.wg.Done()
	
	frameIndex := 0
	ticker := time.NewTicker(s.delay)
	defer ticker.Stop()
	
	for {
		select {
		case <-s.stop:
			return
		case <-ticker.C:
			s.mu.Lock()
			frame := s.frames[frameIndex]
			msg := s.message
			s.mu.Unlock()
			
			fmt.Fprintf(s.writer, "\r%s%s%s %s", ColorCyan, frame, ColorReset, msg)
			
			frameIndex = (frameIndex + 1) % len(s.frames)
		}
	}
}

// WithFrames sets custom spinner frames
func (s *Spinner) WithFrames(frames []string) *Spinner {
	s.frames = frames
	return s
}

// WithDelay sets the animation delay
func (s *Spinner) WithDelay(delay time.Duration) *Spinner {
	s.delay = delay
	return s
}

// WithWriter sets the output writer
func (s *Spinner) WithWriter(w io.Writer) *Spinner {
	s.writer = w
	return s
}

// Simple spinner shortcuts

// Spin runs a function with a spinner
func Spin(message string, fn func() error) error {
	spinner := NewSpinner(message)
	spinner.Start()
	
	err := fn()
	
	if err != nil {
		spinner.Error(err.Error())
		return err
	}
	
	spinner.Success("Done")
	return nil
}

// SpinWithMessage runs a function with a spinner and custom success message
func SpinWithMessage(message, successMsg string, fn func() error) error {
	spinner := NewSpinner(message)
	spinner.Start()
	
	err := fn()
	
	if err != nil {
		spinner.Error(err.Error())
		return err
	}
	
	spinner.Success(successMsg)
	return nil
}

