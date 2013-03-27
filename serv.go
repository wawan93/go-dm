package main

import (
	"fmt"
	"html/template"
	"net/http"
	"strconv"
	"strings"
)

var (
	A, B, C, U, I, D, SD, tmp []int
	mass                      string
	mas                       []int
	a                         int
	err                       error
	f                         bool
)

var cache_templ = template.Must(template.ParseFiles("index.html"))

func check_contain(el int, mas *[]int) bool {
	var i int
	for i = 0; i < len(*mas); i++ {
		if (*mas)[i] == el {
			return true
		}
	}
	return false
}

func union(x, y *[]int) []int {
	var l, i, k int
	l = len(*x) + len(*y) + 1
	var un = make([]int, l)

	copy(un, *x)
	k = len(*x)

	for i = 0; i < len(*y); i++ {
		if !check_contain((*y)[i], &un) {
			k++
			un[k-1] = (*y)[i]
		}
	}
	return un[:k]
}

func inter(x, y, z *[]int) []int {
	var l, i int
	l = len(*x) + len(*y) + len(*z)
	var in = make([]int, l)
	l = 0

	for i = 0; i < len(*x); i++ {
		if check_contain((*x)[i], y) && check_contain((*x)[i], z) {
			in[l] = (*x)[i]
			l++
		}
	}

	return in[:l]
}

func diff(x, y, z *[]int) []int {
	var l, i int
	l = len(*x)
	var dif = make([]int, l)

	l = 0

	for i = 0; i < len(*x); i++ {
		a = (*x)[i]
		if !check_contain(a, y) && !check_contain(a, z) {
			dif[l] = a
			l++
		}
	}

	return dif[:l]
}

func sym_dif(x, y, z *[]int) []int {
	var tmp1 = diff(x, y, z)
	var tmp2 = diff(z, x, y)
	var tmp3 = diff(y, x, z)
	var tmp4 = union(&tmp1, &tmp2)
	tmp4 = union(&tmp4, &tmp3)

	return tmp4
}

func equal(x, y *[]int) bool {
	f = false
	if len(*x) != len(*y) {
		return false
	}
	for i := 0; i < len(*x); i++ {
		a = (*x)[i]
		if !check_contain(a, y) {
			return false
		}
	}
	return true
}

func inn(x, y *[]int) bool {
	if equal(x, y) {
		return false
	}

	for i := 0; i < len(*x); i++ {
		a = (*x)[i]
		if !check_contain(a, y) {
			return false
		}
	}

	return true
}

func enter(s string, mas *[]int) []int {
	var a, k int
	var ss = strings.Fields(s)

	*mas = make([]int, len(ss))

	for _, tmp := range ss {
		a, err = strconv.Atoi(tmp)
		if err != nil {
			panic(err)
		}
		if !check_contain(a, mas) {
			(*mas)[k] = a
			k++
		}
	}

	return *mas
}

func handler(w http.ResponseWriter, r *http.Request) {
	err := cache_templ.ExecuteTemplate(w, "index.html", mass)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func go_handler(w http.ResponseWriter, r *http.Request) {
	A = enter(r.FormValue("A"), &A)
	B = enter(r.FormValue("B"), &B)
	C = enter(r.FormValue("C"), &C)

	fmt.Fprintln(w, " ")
	fmt.Fprintln(w, "Объединение множеств: ")
	tmp = union(&B, &C)
	U = union(&A, &tmp)
	fmt.Fprintln(w, U)

	fmt.Fprintln(w, "Пересечение множеств: ")
	I = inter(&A, &B, &C)
	if len(I) == 0 {
		fmt.Fprintln(w, "пустое множество")
	} else {
		fmt.Fprintln(w, I)
	}

	fmt.Fprintln(w, "Разность A-B-C")
	D = diff(&A, &B, &C)
	if len(D) == 0 {
		fmt.Fprintln(w, "пустое множество")
	} else {
		fmt.Fprintln(w, D)
	}

	fmt.Fprintln(w, "Разность B-A-C")
	D = diff(&B, &A, &C)
	if len(D) == 0 {
		fmt.Fprintln(w, "пустое множество")
	} else {
		fmt.Fprintln(w, D)
	}

	fmt.Fprintln(w, "Разность C-A-B")
	D = diff(&C, &A, &B)
	if len(D) == 0 {
		fmt.Fprintln(w, "пустое множество")
	} else {
		fmt.Fprintln(w, D)
	}

	fmt.Fprintln(w, "Симметрическая разность")
	SD = sym_dif(&A, &B, &C)
	SD = union(&SD, &I)
	fmt.Fprintln(w, SD)

	fmt.Fprintln(w, "<br><br>")

	if equal(&A, &B) {
		fmt.Fprintln(w, "Множествa A и В равны")
	} else {
		fmt.Fprintln(w, "Множествa A и В не равны")
	}

	if inn(&A, &B) {
		fmt.Fprintln(w, "Множество A входит в B")
	} else if inn(&B, &A) {
		fmt.Fprintln(w, "Множество В входит в А")
	} else {
		fmt.Fprintln(w, "Множество А не входит в В")
		fmt.Fprintln(w, "Множество В не входит в А")
	}
}

func main() {
	http.HandleFunc("/", handler)
	http.HandleFunc("/go/", go_handler)
	http.ListenAndServe(":80", nil)
}
