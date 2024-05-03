package hw03frequencyanalysis

import (
	"testing"

	"github.com/stretchr/testify/require"
)

// Change to true if needed.
var taskWithAsteriskIsCompleted = isAsteriksTask

var text = `Как видите, он  спускается  по  лестнице  вслед  за  своим
	другом   Кристофером   Робином,   головой   вниз,  пересчитывая
	ступеньки собственным затылком:  бум-бум-бум.  Другого  способа
	сходить  с  лестницы  он  пока  не  знает.  Иногда ему, правда,
		кажется, что можно бы найти какой-то другой способ, если бы  он
	только   мог   на  минутку  перестать  бумкать  и  как  следует
	сосредоточиться. Но увы - сосредоточиться-то ему и некогда.
		Как бы то ни было, вот он уже спустился  и  готов  с  вами
	познакомиться.
	- Винни-Пух. Очень приятно!
		Вас,  вероятно,  удивляет, почему его так странно зовут, а
	если вы знаете английский, то вы удивитесь еще больше.
		Это необыкновенное имя подарил ему Кристофер  Робин.  Надо
	вам  сказать,  что  когда-то Кристофер Робин был знаком с одним
	лебедем на пруду, которого он звал Пухом. Для лебедя  это  было
	очень   подходящее  имя,  потому  что  если  ты  зовешь  лебедя
	громко: "Пу-ух! Пу-ух!"- а он  не  откликается,  то  ты  всегда
	можешь  сделать вид, что ты просто понарошку стрелял; а если ты
	звал его тихо, то все подумают, что ты  просто  подул  себе  на
	нос.  Лебедь  потом  куда-то делся, а имя осталось, и Кристофер
	Робин решил отдать его своему медвежонку, чтобы оно не  пропало
	зря.
		А  Винни - так звали самую лучшую, самую добрую медведицу
	в  зоологическом  саду,  которую  очень-очень  любил  Кристофер
	Робин.  А  она  очень-очень  любила  его. Ее ли назвали Винни в
	честь Пуха, или Пуха назвали в ее честь - теперь уже никто  не
	знает,  даже папа Кристофера Робина. Когда-то он знал, а теперь
	забыл.
		Словом, теперь мишку зовут Винни-Пух, и вы знаете почему.
		Иногда Винни-Пух любит вечерком во что-нибудь поиграть,  а
	иногда,  особенно  когда  папа  дома,  он больше любит тихонько
	посидеть у огня и послушать какую-нибудь интересную сказку.
		В этот вечер...`

func TestTop10(t *testing.T) {
	t.Run("no words in empty string", func(t *testing.T) {
		require.Len(t, Top10(""), 0)
	})

	t.Run("positive test", func(t *testing.T) {
		if taskWithAsteriskIsCompleted {
			expected := []string{
				"а",         // 8
				"он",        // 8
				"и",         // 6
				"ты",        // 5
				"что",       // 5
				"в",         // 4
				"его",       // 4
				"если",      // 4
				"кристофер", // 4
				"не",        // 4
			}
			require.Equal(t, expected, Top10(text))
		} else {
			expected := []string{
				"он",        // 8
				"а",         // 6
				"и",         // 6
				"ты",        // 5
				"что",       // 5
				"-",         // 4
				"Кристофер", // 4
				"если",      // 4
				"не",        // 4
				"то",        // 4
			}
			require.Equal(t, expected, Top10(text))
		}
	})
}

func TestSplitWords(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		expected    []string
		expectedReg []string
	}{
		{name: "empty", input: "", expected: []string{}, expectedReg: []string{}},
		{
			name:  "en.dog_cat",
			input: "dog,    cat; dog,,,cat	dog...cat ,dog - cat",
			expected: []string{
				"dog,",
				"cat;",
				"dog,,,cat",
				"dog...cat",
				",dog",
				"-",
				"cat",
			},
			expectedReg: []string{
				"dog",
				"cat",
				"dog,,,cat",
				"dog...cat",
				"dog",
				"cat",
			},
		},
		{
			name: "ru.text",
			input: `Предложения  	складываются в абзацы -
			и вы...мы наслаждаетесь	каким-то	очередным ------ бредошедевром?`,
			expected: []string{
				"Предложения",
				"складываются",
				"в",
				"абзацы",
				"-", "и",
				"вы...мы",
				"наслаждаетесь",
				"каким-то",
				"очередным",
				"------",
				"бредошедевром?",
			},
			expectedReg: []string{
				"Предложения",
				"складываются",
				"в",
				"абзацы",
				"и",
				"вы...мы",
				"наслаждаетесь",
				"каким-то",
				"очередным",
				"------",
				"бредошедевром",
			},
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			var result []string
			if taskWithAsteriskIsCompleted {
				result = splitWords(tc.input, pattern)
				require.Equal(t, tc.expectedReg, result)
			} else {
				result = splitWords(tc.input, nil)
				require.Equal(t, tc.expected, result)
			}
		})
	}
}

func TestWordsWidthsSort(t *testing.T) {
	tests := []struct {
		name               string
		input              []string
		expected           []wordWidth
		expectedIgnoreCase []wordWidth
	}{
		{
			name:               "empty",
			input:              []string{},
			expected:           []wordWidth{},
			expectedIgnoreCase: []wordWidth{},
		},
		{
			name:  "en",
			input: []string{"alfa", "beta", "gamma", "Beta", "Alfa"},
			expected: []wordWidth{
				{Word: "Alfa", Width: 1},
				{Word: "Beta", Width: 1},
				{Word: "alfa", Width: 1},
				{Word: "beta", Width: 1},
				{Word: "gamma", Width: 1},
			},
			expectedIgnoreCase: []wordWidth{
				{Word: "alfa", Width: 2},
				{Word: "beta", Width: 2},
				{Word: "gamma", Width: 1},
			},
		},
		{
			name:  "ru",
			input: []string{"Мама", "мыла", "раму,", "раму", "мыла", "мама", "Мыла", "Раму", "мамА"},
			expected: []wordWidth{
				{"мыла", 2},
				{"Мама", 1},
				{"Мыла", 1},
				{"Раму", 1},
				{"мамА", 1},
				{"мама", 1},
				{"раму", 1},
				{"раму,", 1},
			},
			expectedIgnoreCase: []wordWidth{
				{"мама", 3},
				{"мыла", 3},
				{"раму", 2},
				{"раму,", 1},
			},
		},
		// {
		// 	name:          "ru: three_three_two_one",
		// 	isNotCaseSens: true,
		// 	input:         []string{"Мама", "мыла", "раму,", "раму", "мыла", "мама", "Мыла", "Раму", "мамА"},
		// 	expected: []wordWidth{
		// 		{"Мама", 3},
		// 		{"мыла", 3},
		// 		{"раму", 2},
		// 		{"раму,", 1},
		// 	},
		// },
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			result := wordsWidthsSort(tc.input, taskWithAsteriskIsCompleted)
			if taskWithAsteriskIsCompleted {
				require.Equal(t, tc.expectedIgnoreCase, result)
			} else {
				require.Equal(t, tc.expected, result)
			}
		})
	}
}
