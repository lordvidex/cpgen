package main

import (
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/lordvidex/cpgen"
	"github.com/lordvidex/cpgen/pkg/cond"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/progress"
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// the possible pages
const (
	options = iota // the pages to select the options or flags (in CLI mode)
	files          // the pages to edit and add filenames to generate
	loader         // showing update from the file generation channel
)

const (
	padding  = 2
	maxWidth = 80
	eps      = 0.000001
)

var (
	pages []tea.Model
	size  tea.WindowSizeMsg // size of the window
)

// styles
var (
	headerStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#00ff00"))
	errorStyle  = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#ffffff")).
			Background(lipgloss.Color("#dd0000")).
			Padding(1, 4)
	cursorStyle = lipgloss.NewStyle().
			Background(lipgloss.Color("#00ff00")).
			Foreground(lipgloss.Color("#ffffff"))
	normalStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color(lipgloss.Color("#ffffff")))
	docStyle          = lipgloss.NewStyle().Margin(1, 2)
	itemStyle         = lipgloss.NewStyle().PaddingLeft(4)
	selectedItemStyle = lipgloss.NewStyle().PaddingLeft(2).Foreground(lipgloss.Color("170"))
)

// the possible options to generate with
type option struct {
	message string
}

type optionsModel struct {
	selected map[int]struct{}
	help     help.Model
	options  []option
	cursor   int
}

func newOptionsModel(exts []option) *optionsModel {
	return &optionsModel{
		options:  exts,
		selected: make(map[int]struct{}),
		help:     help.New(),
	}
}

// Init implements tea.Model
func (*optionsModel) Init() tea.Cmd {
	return nil
}

// Update implements tea.Model
func (m *optionsModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		size = msg
	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
			return pages[files].Update(size)
		case "ctrl+c", "esc":
			return m, tea.Quit
		case " ":
			_, ok := m.selected[m.cursor]
			if ok {
				delete(m.selected, m.cursor)
			} else {
				m.selected[m.cursor] = struct{}{}
			}
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			} else {
				m.cursor = len(m.options) - 1
			}
		case "down", "j":
			if m.cursor == len(m.options)-1 {
				m.cursor = 0
			} else {
				m.cursor++
			}
		}
	}
	return m, nil
}

// FullHelp implements help.KeyMap
func (*optionsModel) FullHelp() [][]key.Binding {
	return nil
}

// ShortHelp implements help.KeyMap
func (*optionsModel) ShortHelp() []key.Binding {
	return []key.Binding{
		key.NewBinding(
			key.WithKeys(tea.KeySpace.String()),
			key.WithHelp("⎵", "toggle"),
		),
		key.NewBinding(
			key.WithKeys("up", "k"),
			key.WithHelp("↑/up/k", "up"),
		),
		key.NewBinding(
			key.WithKeys("down", "j"),
			key.WithHelp("↓/down/j", "down"),
		),
		key.NewBinding(
			key.WithKeys("enter"),
			key.WithHelp("⏎", "next"),
		),
		key.NewBinding(
			key.WithKeys("esc", "ctrl+c"),
			key.WithHelp("esc/ctrl+c", "quit"),
		),
	}
}

// View implements tea.Model
func (m *optionsModel) View() string {
	header := headerStyle.Render("Select the extensions to generate your solution source codes: ")
	var s strings.Builder
	for i := 0; i < len(m.options); i++ {
		checked := "[ ]"
		if _, ok := m.selected[i]; ok {
			checked = "[✅]"
		}
		row := fmt.Sprintf("%s %s", checked, m.options[i].message)
		if m.cursor == i {
			row = cursorStyle.Render(row)
		} else {
			row = normalStyle.Render(row)
		}
		s.WriteString(row + "\n")
	}
	list := s.String()
	return fmt.Sprintf("%s\n\n%s\n\n%s", header, list, m.help.View(m))
}

// Page 2

type file string

type fileDelegate struct{}

func (d fileDelegate) Height() int                             { return 1 }
func (d fileDelegate) Spacing() int                            { return 0 }
func (d fileDelegate) Update(_ tea.Msg, _ *list.Model) tea.Cmd { return nil }
func (d fileDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	i, ok := listItem.(file)
	if !ok {
		return
	}
	str := fmt.Sprintf("%d. %s", index+1, i)

	if index == m.Index() {
		str = selectedItemStyle.Render("> " + str)
	} else {
		str = itemStyle.Render(str)
	}
	fmt.Fprint(w, str)
}
func (f file) Title() string { return string(f) }

func (f file) Description() string { return string(f) }
func (f file) FilterValue() string {
	return string(f)
}

type filesModel struct {
	l          list.Model
	ti         textinput.Model
	tiHelp     help.Model
	isEditMode bool
}

// FullHelp implements help.KeyMap
func (filesModel) FullHelp() [][]key.Binding {
	return nil
}

// ShortHelp implements help.KeyMap
func (filesModel) ShortHelp() []key.Binding {
	return []key.Binding{
		key.NewBinding(
			key.WithKeys("enter"),
			key.WithHelp("⏎", "save"),
		),
		key.NewBinding(
			key.WithKeys("esc"),
			key.WithHelp("esc", "back"),
		),
	}
}

// Init implements tea.Model
func (f *filesModel) Init() tea.Cmd {
	return nil
}

func (f *filesModel) addItem(item string) tea.Cmd {
	return f.l.InsertItem(len(f.l.Items()), file(item))
}

func (f *filesModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if f.isEditMode {
			switch msg.String() {
			case "enter": // add item and leave edit mode
				cmd := f.addItem(f.ti.Value())
				f.isEditMode = false
				f.ti.SetValue("")
				return f, cmd
			case "esc": // leave edit mode ignoring what was written
				f.isEditMode = false
				f.ti.SetValue("")
				return f, nil
			default: // type text mode
				var cmd tea.Cmd
				f.ti, cmd = f.ti.Update(msg)
				return f, cmd
			}
		} else {
			switch msg.String() {
			case "ctrl+c":
				return f, tea.Quit
			case "esc":
				return pages[options].Update(msg)
			case "enter":
				return pages[loader].Update(size)
			case "backspace", "delete":
				f.l.RemoveItem(f.l.Index())
			case "n":
				f.isEditMode = true
				cmd := f.ti.Focus()
				return f, cmd
			}
		}
	case tea.WindowSizeMsg:
		size = msg
		h, v := docStyle.GetFrameSize()
		f.l.SetSize(msg.Width-h, msg.Height-v)
	}
	var cmd tea.Cmd
	f.l, cmd = f.l.Update(msg)
	return f, cmd
}

func (f *filesModel) View() string {
	if f.isEditMode {
		err := ""
		if f.ti.Err != nil {
			err = errorStyle.Render(f.ti.Err.Error())
		}
		return fmt.Sprintf("%s\n%s\n\n%s", err, f.ti.View(), f.tiHelp.View(f))
	}
	return f.l.View()
}

func newFilesModel(pages []string) *filesModel {
	pg := make([]list.Item, 0, len(pages))
	for _, p := range pages {
		pg = append(pg, file(p))
	}
	listModel := list.New(pg, fileDelegate{}, size.Width, size.Height)
	listModel.AdditionalShortHelpKeys = func() []key.Binding {
		return []key.Binding{
			key.NewBinding(
				key.WithKeys(tea.KeyDelete.String(), tea.KeyBackspace.String()),
				key.WithHelp("⌫/del", "delete"),
			),
			key.NewBinding(
				key.WithKeys("n", "N"),
				key.WithHelp("n", "new"),
			),
		}
	}
	ti := textinput.New()
	ti.Placeholder = "filename..."
	ti.CharLimit = 156
	ti.Width = 20
	ti.Validate = func(s string) error {
		if strings.ContainsRune(strings.TrimSpace(s), ' ') {
			return errors.New("filename should not contain spaces")
		}
		return nil
	}

	h := help.New()

	return &filesModel{
		l:      listModel,
		ti:     ti,
		tiHelp: h,
	}
}

// Page 3

// loadingModel is responsible for showing the progress of file generation to the user in the terminal
type loadingModel struct {
	progress       progress.Model
	spinner        spinner.Model
	update         float64
	finishedUpdate bool
	quitting       bool
	once           sync.Once
}

func (m *loadingModel) Init() tea.Cmd {
	return tea.Batch(m.tickCmd(), m.spinner.Tick)
}

func (m *loadingModel) Bg() {
	exts := pages[options].(*optionsModel).selected
	contains := func(x int) bool {
		_, ok := exts[x]
		return ok
	}
	toStringArr := func(it []list.Item) []string {
		its := make([]string, len(it))
		for i, x := range it {
			its[i] = string(x.(file))
		}
		return its
	}
	ch := cpgen.Generate(toStringArr(pages[files].(*filesModel).l.Items()), cpgen.Config{
		Pq: contains(0),
		Uf: contains(1),
		Sv: contains(2),
		Cf: contains(3),
		FileIO: cond.If(contains(4),
			&cpgen.IO{Input: "input.txt", Output: "output.txt"},
			nil,
		),
	}, "solution")
	for x := range ch {
		m.update = x
	}
	m.finishedUpdate = true
}

func (m *loadingModel) View() string {
	return fmt.Sprintf("\n\n%s %s \n%s", m.spinner.View(), "Generating files... ", m.progress.View())
}

// tickMsg represents a single tick in the progress bar to check for updates
type tickMsg float64

type sig int

const (
	quitting sig = iota
)

func (m *loadingModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds = []tea.Cmd{m.tickCmd(), m.spinner.Tick}
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.progress.Width = msg.Width - padding*2 - 4
		if m.progress.Width > maxWidth {
			m.progress.Width = maxWidth
		}
		m.once.Do(func() {
			go m.Bg()
		})
	case tea.KeyMsg:
		switch msg.String() {
		case tea.KeyCtrlC.String(), "q":
			return m, tea.Quit
		}
	case sig:
		if msg == quitting {
			return m, tea.Quit
		}
	case tickMsg:
		if m.finishedUpdate {
			if !m.quitting {
				m.quitting = true
				// wait for 1 sec to close
				cmds = append(cmds, tea.Tick(time.Second, func(_ time.Time) tea.Msg {
					return quitting
				}))
			} else {
				return m, nil
			}
		}
		cmd := m.progress.SetPercent(m.update)
		cmds = append(cmds, cmd)
	case progress.FrameMsg:
		progressModel, cmd := m.progress.Update(msg)
		m.progress = progressModel.(progress.Model)
		cmds = append(cmds, cmd)
	case spinner.TickMsg:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		cmds = append(cmds, cmd)
	}
	return m, tea.Batch(cmds...)
}

func newLoadingModel() *loadingModel {
	spn := spinner.New()
	spn.Spinner = spinner.Dot
	spn.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
	return &loadingModel{
		progress: progress.New(progress.WithDefaultGradient()),
		spinner:  spn,
	}
}

func (m *loadingModel) tickCmd() tea.Cmd {
	return tea.Tick(time.Millisecond*200, func(t time.Time) tea.Msg {
		return tickMsg(m.update)
	})
}

func main() {
	f, err := os.OpenFile("log.txt", os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		log.Fatal("failed to open log file")
	}
	log.SetOutput(f)
	extensions := []option{
		{message: "Add priority queue template"},
		{message: "Add union find template"},
		{message: "Add sieve of erathostenes template"},
		{message: "Add testcase loop i.e. `t` testcases"},
		{message: "Use files 'input.txt' and 'output.txt' instead of console"},
	}
	defaultPages := []string{"a", "b", "c", "d", "e", "f"}
	pages = []tea.Model{
		newOptionsModel(extensions),
		newFilesModel(defaultPages),
		newLoadingModel(),
	}
	p := tea.NewProgram(pages[options], tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
}
