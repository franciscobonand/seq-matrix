package entity

import (
	"testing"
)

func TestSequences_Validate(t *testing.T) {
	tests := []struct {
		name    string
		letters []string
		want    bool
		wantErr bool
	}{
		{
			name:    "Valid Sequece - Column, Row and Secondary Diagonal",
			letters: []string{"DUHBHB", "DUBUHD", "UBUUHU", "BHBDHH", "DDDDUB", "UDBDUH"},
			want:    true,
			wantErr: false,
		},
		{
			name:    "Valid Sequece - Primary and Secondary Diagonals",
			letters: []string{"BUHDHB", "DBHDHD", "UUBUHU", "BHUBUH", "HDHUDB", "UDUDUH"},
			want:    true,
			wantErr: false,
		},
		{
			name:    "Invalid Sequece - No repetitions",
			letters: []string{"BUHDHB", "DBHUHD", "UUBUUU", "BHBDHH", "HDHUDB", "UDBDUH"},
			want:    false,
			wantErr: false,
		},
		{
			name:    "Invalid Sequece - One repetition",
			letters: []string{"BUHDHB", "DBHHHD", "UUHUUU", "BHBDHH", "HDHUDB", "UDBDUH"},
			want:    false,
			wantErr: false,
		},
		{
			name:    "Invalid Input - Length",
			letters: []string{"DUHBHB", "UBUUHU", "BHBDHH", "DDDDUB", "UDBDUH"},
			want:    false,
			wantErr: true,
		},
		{
			name:    "Invalid Input - Wrong alphabet",
			letters: []string{"DBHUHD", "RSATPB", "BHBDHH", "HDHUDB", "UDBDUH"},
			want:    false,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := Sequences{
				Letters: tt.letters,
			}
			got, err := s.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Sequences.Validate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Sequences.Validate() = %v, want %v", got, tt.want)
			}
		})
	}
}
