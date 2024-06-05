package iterator

import (
	"fmt"
)

type Aggragate interface {
	interator()
}

type Iterator interface {
	hasNext() bool
	next() Book
}

type Book struct {
	name string
}
func NewBook(name string) Book {
	return Book{name: name}
}
func (b Book) String() string {
	return b.name
}

type BookShelf struct {
	Books []Book
	Last int
}
func NewBookShelf(maxSize int) *BookShelf {
	return &BookShelf{
		Books: make([]Book, maxSize),
	}
}
func (bs *BookShelf) getBookAt(index int) Book{
	return bs.Books[index]
}
func (bs *BookShelf) appendBook(book Book) {
	// if bs.Last >= len(bs.Books) {
	// 	fmt.Println("Bookshelf if full!")
	// 	return
	// }
	bs.Books = append(bs.Books, book)
	bs.Last++
}
func (bs *BookShelf) getLength() int {
	return bs.Last
}
func (bs *BookShelf) iterator() Iterator{
	return NewBookShelfIterator(bs)
}

type BookShelfIterator struct {
	BookShelf *BookShelf
	Index int
}
func NewBookShelfIterator(bs *BookShelf) *BookShelfIterator {
	return &BookShelfIterator{BookShelf: bs}
}
func (bsi *BookShelfIterator) hasNext() bool {
	fmt.Println(bsi.Index, bsi.BookShelf.getLength())
	return bsi.Index <= bsi.BookShelf.getLength()
}
func (bsi *BookShelfIterator) next() Book{
	book := bsi.BookShelf.getBookAt(bsi.Index)
	bsi.Index++
	return book
}

func Exec() {
	shelf := NewBookShelf(1)
	shelf.appendBook(NewBook("Book 1"))
	shelf.appendBook(NewBook("Book 2"))
	shelf.appendBook(NewBook("Book 3"))
	shelf.appendBook(NewBook("Book 4"))

	internal := shelf.iterator()

	for internal.hasNext() {
		book := internal.next()
		fmt.Println(book.String())
	}
}
