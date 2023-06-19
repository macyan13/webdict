package query

//
// func TestAllUsersHandler_Handle(t *testing.T) {
//	type fields struct {
//		userRepo UserViewRepository
//	}
//	tests := []struct {
//		name     string
//		fieldsFn func() fields
//		want     []UserView
//		wantErr  assert.ErrorAssertionFunc
//	}{
//		{
//			"Error on getting users from DB",
//			func() fields {
//				userRepo := MockUserViewRepository{}
//				userRepo.On("GetAllViews").Return(nil, errors.New("testErr"))
//				return fields{userRepo: &userRepo}
//			},
//			nil,
//			func(t assert.TestingT, err error, i ...interface{}) bool {
//				assert.Equal(t, "testErr", err.Error(), i)
//				return true
//			},
//		},
//		{
//			"Sanitize is called",
//			func() fields {
//				userRepo := MockUserViewRepository{}
//				users := []UserView{{
//					ID:    "testId",
//					Name:  `<a href="javascript:alert('XSS1')" onmouseover="alert('XSS2')">TestName<a>`,
//					Email: `<a href="javascript:alert('XSS1')" onmouseover="alert('XSS2')">TestEmail<a>`,
//					ID:  2,
//				}}
//				userRepo.On("GetAllViews").Return(users, nil)
//				return fields{userRepo: &userRepo}
//			},
//			[]UserView{{
//				ID:    "testId",
//				Name:  "TestName",
//				Email: "TestEmail",
//				ID:  2,
//			}},
//			func(t assert.TestingT, err error, i ...interface{}) bool {
//				assert.Nil(t, err, i)
//				return false
//			},
//		},
//	}
//
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			h := NewAllUsersHandler(tt.fieldsFn().userRepo)
//			got, err := h.Handle()
//			if !tt.wantErr(t, err, "Handle()") {
//				return
//			}
//			assert.Equalf(t, tt.want, got, "Handle()")
//		})
//	}
//}
