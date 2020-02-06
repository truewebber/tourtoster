package user

import (
	"testing"
)

func TestUser_HasPermission(t *testing.T) {
	type fields struct {
		Permissions Permission
	}
	type args struct {
		p Permission
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
	}{
		{
			name: "test1",
			fields: fields{
				Permissions: CreateNewUserPermission | EditToursPermission | EditAllBookingsPermission | EditUserBookingsPermission,
			},
			args: args{p: CreateNewUserPermission},
			want: true,
		},
		{
			name: "test2",
			fields: fields{
				Permissions: CreateNewUserPermission | EditToursPermission | EditAllBookingsPermission | EditUserBookingsPermission,
			},
			args: args{p: EditToursPermission},
			want: true,
		},
		{
			name: "test3",
			fields: fields{
				Permissions: CreateNewUserPermission | EditToursPermission | EditAllBookingsPermission | EditUserBookingsPermission,
			},
			args: args{p: EditAllBookingsPermission},
			want: true,
		},
		{
			name: "test4",
			fields: fields{
				Permissions: CreateNewUserPermission | EditToursPermission | EditAllBookingsPermission | EditUserBookingsPermission,
			},
			args: args{p: EditUserBookingsPermission},
			want: true,
		},
		{
			name: "test5",
			fields: fields{
				Permissions: EditToursPermission | EditAllBookingsPermission | EditUserBookingsPermission,
			},
			args: args{p: CreateNewUserPermission},
			want: false,
		},
		{
			name: "test6",
			fields: fields{
				Permissions: CreateNewUserPermission | EditAllBookingsPermission | EditUserBookingsPermission,
			},
			args: args{p: EditToursPermission},
			want: false,
		},
		{
			name: "test7",
			fields: fields{
				Permissions: CreateNewUserPermission | EditToursPermission | EditUserBookingsPermission,
			},
			args: args{p: EditAllBookingsPermission},
			want: false,
		},
		{
			name: "test8",
			fields: fields{
				Permissions: CreateNewUserPermission | EditToursPermission | EditAllBookingsPermission,
			},
			args: args{p: EditUserBookingsPermission},
			want: false,
		},
		{
			name:   "test9",
			fields: fields{},
			args:   args{p: EditUserBookingsPermission},
			want:   false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := &User{
				Permissions: tt.fields.Permissions,
			}
			if got := u.HasPermission(tt.args.p); got != tt.want {
				t.Errorf("HasPermission() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUser_AdminPage(t *testing.T) {
	type fields struct {
		Permissions Permission
	}
	tests := []struct {
		name   string
		fields fields
		want   bool
	}{
		{
			name:   "test1",
			fields: fields{Permissions: CreateNewUserPermission | EditToursPermission | EditAllBookingsPermission | EditUserBookingsPermission},
			want:   true,
		},
		{
			name:   "test2",
			fields: fields{Permissions: CreateNewUserPermission | EditToursPermission | EditAllBookingsPermission},
			want:   true,
		},
		{
			name:   "test3",
			fields: fields{Permissions: CreateNewUserPermission | EditToursPermission},
			want:   true,
		},
		{
			name:   "test4",
			fields: fields{Permissions: CreateNewUserPermission},
			want:   true,
		},
		{
			name:   "test5",
			fields: fields{},
			want:   false,
		},
		{
			name:   "test6",
			fields: fields{Permissions: EditUserBookingsPermission},
			want:   false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := &User{
				Permissions: tt.fields.Permissions,
			}
			if got := u.AdminPage(); got != tt.want {
				t.Errorf("AdminPage() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestShortName(t *testing.T) {
	type args struct {
		u *User
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "test1",
			args: args{u: &User{
				FirstName: "Alex",
			}},
			want: "A.",
		},
		{
			name: "test2",
			args: args{u: &User{
				FirstName:  "Alex",
				SecondName: "Moshi",
			}},
			want: "A.M.",
		},
		{
			name: "test3",
			args: args{u: &User{
				FirstName:  "anna",
				SecondName: "maria",
			}},
			want: "a.m.",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ShortName(tt.args.u); got != tt.want {
				t.Errorf("ShortName() = %v, want %v", got, tt.want)
			}
		})
	}
}
