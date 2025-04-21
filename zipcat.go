package zipcat

import (
	"archive/zip"
	"bufio"
	"errors"
	"io"
	"iter"
	"os"
)

type ZipReaderLike struct {
	io.ReaderAt
	Size int64
}

func (l ZipReaderLike) ToZipReader() (*zip.Reader, error) {
	return zip.NewReader(l.ReaderAt, l.Size)
}

type ZipWriterLike struct {
	io.Writer
}

func (l ZipWriterLike) ToZipWriter() *zip.Writer {
	return zip.NewWriter(l.Writer)
}

type ZipWriter struct{ *zip.Writer }

func (w ZipWriter) Close() error { return w.Writer.Close() }

func (w ZipWriter) AppendFiles(files []*zip.File) error {
	for _, zfile := range files {
		e := w.Writer.Copy(zfile)
		if nil != e {
			return e
		}
	}
	return nil
}

func (w ZipWriter) AppendZipLike(l ZipReaderLike) error {
	rdr, e := l.ToZipReader()
	if nil != e {
		return e
	}

	var files []*zip.File = rdr.File
	return w.AppendFiles(files)
}

func (w ZipWriter) AppendZipFile(f *os.File) error {
	finfo, e := f.Stat()
	if nil != e {
		return e
	}
	var size int64 = finfo.Size()
	zlike := ZipReaderLike{
		ReaderAt: f,
		Size:     size,
	}
	return w.AppendZipLike(zlike)
}

func (w ZipWriter) AppendZip(filename string) error {
	f, e := os.Open(filename)
	if nil != e {
		return e
	}
	defer f.Close()
	return w.AppendZipFile(f)
}

func (w ZipWriter) AppendZipFiles(filenames iter.Seq[string]) error {
	for fname := range filenames {
		e := w.AppendZip(fname)
		if nil != e {
			return e
		}
	}
	return nil
}

func (l ZipWriterLike) ConcatenateZipfiles(filenames iter.Seq[string]) error {
	wtr := ZipWriter{Writer: l.ToZipWriter()}
	return errors.Join(
		wtr.AppendZipFiles(filenames),
		wtr.Close(),
	)
}

func ConcatenatedZipfilesToStdout(filenames iter.Seq[string]) error {
	wlike := ZipWriterLike{Writer: os.Stdout}
	return wlike.ConcatenateZipfiles(filenames)
}

func StdinToZipFilenamesToConcatenatedToStdout() error {
	var s *bufio.Scanner = bufio.NewScanner(os.Stdin)
	return ConcatenatedZipfilesToStdout(func(
		yield func(string) bool,
	) {
		for s.Scan() {
			var filename string = s.Text()
			if !yield(filename) {
				return
			}
		}
	})
}
