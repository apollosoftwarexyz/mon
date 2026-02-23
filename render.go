package mon

import (
	"fmt"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/apollosoftwarexyz/mon/animations"
	"github.com/apollosoftwarexyz/mon/formatting"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	boldStyle     = lipgloss.NewStyle().Bold(true)
	completeStyle = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("34"))
	errorStyle    = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("160"))
)

const (
	completeIcon = "✓"
	errorIcon    = "✖"
)

type notifyMsg struct {
	tag int
}

type doneMsg struct{}

type model struct {
	prog         *tea.Program
	spinnerAnim  *animations.A
	ellipsisAnim *animations.A
	caption      string
	start        time.Time
	done         bool

	notifyMutex sync.Mutex

	tasksMutex sync.RWMutex
	tasks      []Task
}

func (m *model) notify() {
	// Use a mutex around the notifyMsg to ensure we do not queue simultaneous
	// notify messages (we should drop them instead).
	if m.notifyMutex.TryLock() {
		defer m.notifyMutex.Unlock()
		m.prog.Send(notifyMsg{})
	}
}

func (m *model) addTask(task Task) {
	m.tasksMutex.Lock()
	defer m.tasksMutex.Unlock()
	m.tasks = append(m.tasks, task)
}

func (m *model) Init() tea.Cmd {
	return nil
}

func (m *model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case notifyMsg:
		return m, nil
	case doneMsg:
		m.done = true
		return m, tea.Quit
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			m.done = true
			return m, tea.Quit
		}
	}

	return m, nil
}

func (m *model) View() string {
	if m.done {
		return ""
	}

	var s strings.Builder

	t := time.Since(m.start)
	spinner := m.spinnerAnim.Frame(t)

	{
		// Ensure the tasks are not currently being managed.
		m.tasksMutex.RLock()
		defer m.tasksMutex.RUnlock()

		inProgressTasks := make([]Task, 0)

		for _, t := range m.tasks {
			if t.IsError() {
				if time.Since(t.GetCompletedAt()).Seconds() > 15 {
					continue
				}
			} else if t.IsCompleted() {
				if time.Since(t.GetCompletedAt()).Seconds() > 2 {
					continue
				}
			}

			inProgressTasks = append(inProgressTasks, t)
		}

		for _, t := range inProgressTasks {
			s.WriteString(m.renderTask(t, inProgressTasks, spinner))
		}
	}

	s.WriteRune('\n')

	s.WriteString(boldStyle.Render(fmt.Sprintf("%s (%0.1fs) %s%s\n", spinner, float64(t.Milliseconds())/1000, m.caption, m.ellipsisAnim.Frame(t))))

	return s.String()
}

func getLongestNameLength(tasks []Task) int {
	l := 0

	for _, t := range tasks {
		nameLength := len(t.GetName())
		if nameLength > l {
			l = nameLength
		}
	}

	return l
}

func renderProgress(t Task) string {
	return t.GetUnit().RenderProgress(t.GetCompleteSteps(), t.GetTotalSteps())
}

func getLongestProgressLength(allTasks []Task) int {
	l := 0

	for _, t := range allTasks {
		if formattedLen := len(renderProgress(t)); formattedLen > l {
			l = formattedLen
		}
	}

	return l
}

func (m *model) renderTask(t Task, allTasks []Task, spinner string) string {
	var s strings.Builder

	icon := spinner
	if t.IsCompleted() {
		if t.IsError() {
			icon = errorIcon
		} else {
			icon = completeStyle.Render(completeIcon)
		}
	}
	s.WriteString(icon)
	s.WriteRune(' ')
	if name := t.GetName(); name != "" {
		s.WriteString(fmt.Sprintf("%"+strconv.Itoa(getLongestNameLength(allTasks))+"s", name))

		if t.GetCaption() != "" {
			s.WriteString(": ")
		} else {
			s.WriteString(" ")
		}
	}

	if t.GetCaption() != "" {
		s.WriteString(fmt.Sprintf("%-20s", t.GetCaption()))
		s.WriteRune(' ')
	}

	s.WriteString(fmt.Sprintf("| %s", formatting.Duration(t.GetElapsed())))
	s.WriteString(" ")

	if t.IsError() {
		s.WriteString("| ")
		s.WriteString(t.GetError().Error())
	}

	if !t.IsIndeterminate() {
		s.WriteString("| ")
		s.WriteString(fmt.Sprintf("%"+strconv.Itoa(getLongestProgressLength(allTasks))+"s", renderProgress(t)))
		s.WriteString(" ")
	}

	estimatedCompletion, hasEstimatedCompletion := t.GetEstimatedCompletion()
	if hasEstimatedCompletion {
		s.WriteString("| ")
		s.WriteString(fmt.Sprintf("eta: %5s |", formatting.Duration(estimatedCompletion)))
	}

	if !t.IsCompleted() {
		averageTimePerStep, hasAverageTimePerStep := t.GetAverageTimePerStep()
		if hasAverageTimePerStep {
			s.WriteRune(' ')
			s.WriteString(fmt.Sprintf("%s/s", t.GetUnit().Render(uint64(1/averageTimePerStep.Seconds()))))
		}
	}

	if t.IsError() {
		return errorStyle.Render(s.String()) + "\n"
	}

	s.WriteRune('\n')
	return s.String()
}
