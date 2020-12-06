package storage

import (
	"io"
	"path"
	"reflect"
	"strings"
	"testing"

	"github.com/spf13/afero"
)

/*
   TESTS
*/
func TestNewFileStore(t *testing.T) {
	fs = afero.NewMemMapFs()
	type args struct {
		path string
	}
	tests := []struct {
		name    string
		args    args
		want    *FileStore
		wantErr bool
	}{
		{
			name:    "Default",
			args:    args{path: "/tmp/media"},
			want:    &FileStore{Path: "/tmp/media"},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetFileStore(tt.args.path)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetFileStore() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetFileStore() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFileStore_Save(t *testing.T) {
	fs = afero.NewMemMapFs()
	storePath := afero.GetTempDir(fs, "media")

	type fields struct {
		Path string
	}
	type args struct {
		file io.ReadSeeker
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    string
		wantErr bool
	}{
		{
			name:    "NonExists",
			fields:  fields{Path: storePath},
			args:    args{file: strings.NewReader("  ?HelloWorldTest!  ")},
			want:    "1fe569ab5a74d6bf7c7a783fcc61dfc30cba304628e31547c19135dd24f040d5",
			wantErr: false,
		},
		{
			name:    "Exists",
			fields:  fields{Path: storePath},
			args:    args{file: strings.NewReader("  Test  ")},
			want:    "966d77a20be11045ac1ffa0f42f8a97569e8ba70966b287575899d875bf62b9e",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := FileStore{
				Path: tt.fields.Path,
			}
			got, err := s.SaveFile(tt.args.file)
			if (err != nil) != tt.wantErr {
				t.Errorf("SaveFile() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("SaveFile() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFileStore_Get(t *testing.T) {
	fs = afero.NewMemMapFs()
	storePath := afero.GetTempDir(fs, "media")
	err := afero.WriteFile(fs, path.Join(storePath, "testfile"), []byte("  ?HelloWorldTest!  "), 0755)
	if err != nil {
		t.Errorf("GetFileWithHash() error in setup: %v", err)
	}
	wantReader, err := fs.Open(path.Join(storePath, "testfile"))
	if err != nil {
		t.Errorf("GetFileWithHash() error in setup: %v", err)
	}

	type fields struct {
		Path string
	}
	type args struct {
		hash string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    io.ReadSeeker
		wantErr bool
	}{
		{
			name:    "Exists",
			fields:  fields{Path: storePath},
			args:    args{hash: "testfile"},
			want:    wantReader,
			wantErr: false,
		},
		{
			name:    "NonExists",
			fields:  fields{Path: storePath},
			args:    args{hash: "testfile2"},
			want:    nil,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := FileStore{
				Path: tt.fields.Path,
			}
			got, err := s.GetFileWithHash(tt.args.hash)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetFileWithHash() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetFileWithHash() got = %v, want %v", got, tt.want)
			}
		})
	}
}

/*
   BENCHMARKS
*/

func BenchmarkFileStore_SaveFile(b *testing.B) {
	for i := 0; i < b.N; i++ {
		b.StopTimer()
		reader := strings.NewReader("HelloWorld")
		filestore := setupFileStoreWithTempDir()
		b.StartTimer()

		_, err := filestore.SaveFile(reader)
		if err != nil {
			b.Errorf("SaveFile() error = %v, wantErr nil", err)
		}
	}
}

func setupFileStoreWithTempDir() FileStore {
	fs = afero.NewMemMapFs()
	storePath := afero.GetTempDir(fs, "media")
	return FileStore{Path: storePath}
}
