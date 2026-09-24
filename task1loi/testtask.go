package main

import (
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

type Grammar struct {
	N     []string
	Sigma []string
	P     map[string][]string
	S     string
}

func (g *Grammar) isNonterminal(s string) bool {
	for _, n := range g.N {
		if n == s {
			return true
		}
	}
	return false
}

func (g *Grammar) startsWith(chain, prefix string) (bool, string) {
	if strings.HasPrefix(chain, prefix) {
		return true, chain[len(prefix):]
	}
	return false, ""
}

func (g *Grammar) generateNewNonterminal(base string) string {
	name := base + "'"
	for g.isNonterminal(name) {
		name += "'"
	}
	return name
}

func (g *Grammar) eliminateDirectLeftRecursion(A string) {
	rules := g.P[A]
	var alpha []string
	var beta []string

	for _, rule := range rules {
		if strings.HasPrefix(rule, A) {
			alpha = append(alpha, rule[len(A):])
		} else {
			beta = append(beta, rule)
		}
	}

	if len(alpha) == 0 {
		return
	}

	newN := g.generateNewNonterminal(A)

	newRulesA := []string{}
	for _, b := range beta {
		newRulesA = append(newRulesA, b+newN)
	}
	g.P[A] = newRulesA

	newRulesNew := []string{}
	for _, a := range alpha {
		newRulesNew = append(newRulesNew, a+newN)
	}
	newRulesNew = append(newRulesNew, "ε")
	g.P[newN] = newRulesNew

	g.N = append(g.N, newN)
}

func (g *Grammar) eliminateIndirectLeftRecursion(i int) {
	Ai := g.N[i]
	rules := g.P[Ai]
	newRules := []string{}

	for _, rule := range rules {
		substituted := false
		for j := 0; j < i; j++ {
			Aj := g.N[j]
			if ok, rest := g.startsWith(rule, Aj); ok {
				for _, AjRule := range g.P[Aj] {
					newRules = append(newRules, AjRule+rest)
				}
				substituted = true
				break
			}
		}
		if !substituted {
			newRules = append(newRules, rule)
		}
	}
	g.P[Ai] = newRules
}

func (g *Grammar) EliminateLeftRecursion() {
	n := len(g.N)
	for i := 0; i < n; i++ {
		g.eliminateIndirectLeftRecursion(i)
		g.eliminateDirectLeftRecursion(g.N[i])
	}
}

func parseList(input string) []string {
	parts := strings.Split(input, ",")
	result := []string{}
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p != "" {
			result = append(result, p)
		}
	}
	return result
}

func parseRules(input string) map[string][]string {
	rules := make(map[string][]string)
	parts := strings.Split(input, ",")
	for _, part := range parts {
		part = strings.TrimSpace(part)
		if part == "" {
			continue
		}
		idx := strings.Index(part, "->")
		if idx < 0 {
			continue
		}
		left := strings.TrimSpace(part[:idx])
		right := strings.TrimSpace(part[idx+2:])
		alternatives := strings.Split(right, "|")
		for _, alt := range alternatives {
			alt = strings.TrimSpace(alt)
			if alt != "" {
				rules[left] = append(rules[left], alt)
			}
		}
	}
	return rules
}

func (g *Grammar) String() string {
	var sb strings.Builder
	sb.WriteString("N = {" + strings.Join(g.N, ", ") + "}\n")
	sb.WriteString("Σ = {" + strings.Join(g.Sigma, ", ") + "}\n")
	sb.WriteString("S = " + g.S + "\n")
	sb.WriteString("P:\n")
	for _, n := range g.N {
		if rules, ok := g.P[n]; ok {
			sb.WriteString("  " + n + " -> " + strings.Join(rules, " | ") + "\n")
		}
	}
	return sb.String()
}

func main() {
	a := app.New()
	w := a.NewWindow("Удаление левой рекурсии в КС-грамматике")
	w.Resize(fyne.NewSize(600, 600))

	entryN := widget.NewEntry()
	entryN.SetPlaceHolder("Нетерминалы через запятую: S, A, B")

	entrySigma := widget.NewEntry()
	entrySigma.SetPlaceHolder("Терминалы через запятую: a, b")

	entryP := widget.NewEntry()
	entryP.SetPlaceHolder("Правила: S -> Ba | Ab, A -> Sa | AAb | a, B -> Sb | BBa | b")

	entryS := widget.NewEntry()
	entryS.SetPlaceHolder("Начальный символ: S")

	labelRes := widget.NewLabel("Результат появится здесь...")
	labelRes.Wrapping = fyne.TextWrapWord

	btn := widget.NewButton("Удалить левую рекурсию", func() {
		g := &Grammar{}
		g.N = parseList(entryN.Text)
		g.Sigma = parseList(entrySigma.Text)
		g.P = parseRules(entryP.Text)
		g.S = strings.TrimSpace(entryS.Text)

		if len(g.N) == 0 {
			labelRes.SetText("Ошибка: не заданы нетерминалы")
			return
		}
		if g.S == "" {
			labelRes.SetText("Ошибка: не задан начальный символ")
			return
		}

		g.EliminateLeftRecursion()

		labelRes.SetText(g.String())
	})

	btnExample := widget.NewButton("Загрузить пример (задание 5.1)", func() {
		entryN.SetText("S, A, B")
		entrySigma.SetText("a, b")
		entryP.SetText("S -> Ba | Ab, A -> Sa | AAb | a, B -> Sb | BBa | b")
		entryS.SetText("S")
	})

	w.SetContent(container.NewVBox(
		widget.NewLabel("Нетерминалы (N):"),
		entryN,
		widget.NewLabel("Терминалы (Σ):"),
		entrySigma,
		widget.NewLabel("Правила (P):"),
		entryP,
		widget.NewLabel("Начальный символ (S):"),
		entryS,
		btnExample,
		btn,
		widget.NewSeparator(),
		widget.NewLabel("Результат:"),
		labelRes,
	))

	w.ShowAndRun()
}
