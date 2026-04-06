package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// Panel indices
const (
	panelChanges  = 0
	panelBranches = 1
	panelCommits  = 2
)

type model struct {
	changes       []string
	branches      []string
	commits       []string
	currentBranch string

	activePanel int
	cursors     [3]int
	offsets     [3]int

	width  int
	height int
}

// Styles
var (
	titleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("229")).
			Background(lipgloss.Color("57")).
			Padding(0, 1)

	activeBorderStyle = lipgloss.NewStyle().
				Border(lipgloss.RoundedBorder()).
				BorderForeground(lipgloss.Color("57"))

	inactiveBorderStyle = lipgloss.NewStyle().
				Border(lipgloss.RoundedBorder()).
				BorderForeground(lipgloss.Color("240"))

	selectedStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("229")).
			Background(lipgloss.Color("57"))

	cursorStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("212"))

	branchCurrentStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("114"))

	statusAddedStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("114"))

	statusModifiedStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("214"))

	statusDeletedStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("203"))

	hashStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("214"))

	dimStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("240"))
)

func git(args ...string) []string {
	cmd := exec.Command("git", args...)
	out, err := cmd.Output()
	if err != nil {
		return nil
	}
	s := strings.TrimSpace(string(out))
	if s == "" {
		return nil
	}
	return strings.Split(s, "\n")
}

func currentBranch() string {
	lines := git("rev-parse", "--abbrev-ref", "HEAD")
	if len(lines) > 0 {
		return lines[0]
	}
	return ""
}

func loadChanges() []string {
	return git("status", "--porcelain")
}

func loadBranches() []string {
	raw := git("branch", "--format=%(refname:short)")
	return raw
}

func loadCommits(branch string) []string {
	if branch == "" {
		return nil
	}
	return git("log", branch, "--oneline", "-30")
}

func initialModel() model {
	branches := loadBranches()
	cur := currentBranch()
	cursorIdx := 0
	for i, b := range branches {
		if b == cur {
			cursorIdx = i
			break
		}
	}
	m := model{
		changes:       loadChanges(),
		branches:      branches,
		currentBranch: cur,
		activePanel:   panelChanges,
	}
	m.cursors[panelBranches] = cursorIdx
	m.commits = loadCommits(m.selectedBranch())
	return m
}

func (m model) selectedBranch() string {
	if len(m.branches) == 0 {
		return ""
	}
	return m.branches[m.cursors[panelBranches]]
}

func (m model) panelItems(panel int) []string {
	switch panel {
	case panelChanges:
		return m.changes
	case panelBranches:
		return m.branches
	case panelCommits:
		return m.commits
	}
	return nil
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil

	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit

		case "tab":
			m.activePanel = (m.activePanel + 1) % 3
			return m, nil

		case "shift+tab":
			m.activePanel = (m.activePanel + 2) % 3
			return m, nil

		case "j", "down":
			items := m.panelItems(m.activePanel)
			if m.cursors[m.activePanel] < len(items)-1 {
				m.cursors[m.activePanel]++
			}
			if m.activePanel == panelBranches {
				m.commits = loadCommits(m.selectedBranch())
				m.cursors[panelCommits] = 0
				m.offsets[panelCommits] = 0
			}
			return m, nil

		case "k", "up":
			if m.cursors[m.activePanel] > 0 {
				m.cursors[m.activePanel]--
			}
			if m.activePanel == panelBranches {
				m.commits = loadCommits(m.selectedBranch())
				m.cursors[panelCommits] = 0
				m.offsets[panelCommits] = 0
			}
			return m, nil

		case "r":
			m.changes = loadChanges()
			m.branches = loadBranches()
			m.commits = loadCommits(m.selectedBranch())
			return m, nil
		}
	}
	return m, nil
}

func (m model) View() string {
	if m.width == 0 {
		return "Loading..."
	}

	innerWidth := m.width - 4 // border + padding

	changesHeight := m.height / 4
	branchesHeight := m.height / 4
	commitsHeight := m.height - changesHeight - branchesHeight - 6 // titles + borders

	if changesHeight < 3 {
		changesHeight = 3
	}
	if branchesHeight < 3 {
		branchesHeight = 3
	}
	if commitsHeight < 3 {
		commitsHeight = 3
	}

	changesView := m.renderPanel(panelChanges, innerWidth, changesHeight)
	branchesView := m.renderPanel(panelBranches, innerWidth, branchesHeight)
	commitsView := m.renderPanel(panelCommits, innerWidth, commitsHeight)

	// Titles
	changesTitle := titleStyle.Render(fmt.Sprintf(" Changes (%d) ", len(m.changes)))
	branchesTitle := titleStyle.Render(fmt.Sprintf(" Branches (%d) ", len(m.branches)))

	commitLabel := m.selectedBranch()
	if commitLabel == "" {
		commitLabel = "none"
	}
	commitsTitle := titleStyle.Render(fmt.Sprintf(" Commits · %s ", commitLabel))

	borderFn := func(panel int) lipgloss.Style {
		if panel == m.activePanel {
			return activeBorderStyle.Width(innerWidth)
		}
		return inactiveBorderStyle.Width(innerWidth)
	}

	help := dimStyle.Render("  tab: switch panel · j/k: navigate · r: refresh · q: quit")

	return lipgloss.JoinVertical(lipgloss.Left,
		changesTitle,
		borderFn(panelChanges).Render(changesView),
		branchesTitle,
		borderFn(panelBranches).Render(branchesView),
		commitsTitle,
		borderFn(panelCommits).Render(commitsView),
		help,
	)
}

func (m *model) renderPanel(panel, width, height int) string {
	items := m.panelItems(panel)
	cursor := m.cursors[panel]

	// Scroll offset
	if cursor < m.offsets[panel] {
		m.offsets[panel] = cursor
	}
	if cursor >= m.offsets[panel]+height {
		m.offsets[panel] = cursor - height + 1
	}

	var lines []string
	for i := m.offsets[panel]; i < len(items) && i < m.offsets[panel]+height; i++ {
		line := items[i]
		rendered := m.renderLine(panel, i, line, width)
		if i == cursor && panel == m.activePanel {
			rendered = selectedStyle.Width(width).Render(rendered)
		} else if i == cursor {
			rendered = cursorStyle.Width(width).Render(rendered)
		}
		lines = append(lines, rendered)
	}

	// Pad remaining lines
	for len(lines) < height {
		lines = append(lines, strings.Repeat(" ", width))
	}

	return strings.Join(lines, "\n")
}

func (m model) renderLine(panel, idx int, line string, width int) string {
	switch panel {
	case panelChanges:
		if len(line) < 3 {
			return line
		}
		status := line[:2]
		file := strings.TrimSpace(line[2:])
		switch {
		case strings.Contains(status, "A"), strings.Contains(status, "?"):
			return statusAddedStyle.Render("+ ") + file
		case strings.Contains(status, "D"):
			return statusDeletedStyle.Render("- ") + file
		default:
			return statusModifiedStyle.Render("~ ") + file
		}

	case panelBranches:
		if line == m.currentBranch {
			return branchCurrentStyle.Render("● " + line)
		}
		return "  " + line

	case panelCommits:
		parts := strings.SplitN(line, " ", 2)
		if len(parts) == 2 {
			return hashStyle.Render(parts[0]) + " " + parts[1]
		}
		return line
	}
	return line
}

func main() {
	// Check we're in a git repo
	if err := exec.Command("git", "rev-parse", "--git-dir").Run(); err != nil {
		fmt.Fprintln(os.Stderr, "Not a git repository")
		os.Exit(1)
	}

	p := tea.NewProgram(initialModel(), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
