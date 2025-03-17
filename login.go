package main

// import (
//     "fmt"
//     "log"

//     "github.com/go-gl/gl/v3.3-core/gl"
//     "github.com/go-gl/glfw/v3.2/glfw"
// )

// type LoginWindow struct {
//     win          *glfw.Window
//     username    string
//     password    string
//     focusField  int // 0 for username, 1 for password
//     done        bool
//     success     bool
//     errMsg      string
// }

// func NewLoginWindow() (*LoginWindow, error) {
//     // Set window hints for login window
//     glfw.WindowHint(glfw.Resizable, glfw.False)  // Login window shouldn't be resizable
//     glfw.WindowHint(glfw.ContextVersionMajor, 3)
//     glfw.WindowHint(glfw.ContextVersionMinor, 3)
//     glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
//     glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True)

//     win, err := glfw.CreateWindow(400, 300, "Login - Gocraft", nil, nil)
//     if err != nil {
//         return nil, err
//     }

//     login := &LoginWindow{
//         win:         win,
//         username:    "",
//         password:    "",
//         focusField:  0,
//     }

//     win.SetKeyCallback(login.handleKey)
//     win.SetCharCallback(login.handleChar)

//     return login, nil
// }

// func (l *LoginWindow) handleKey(w *glfw.Window, key glfw.Key, scancode int, action glfw.Action, mods glfw.ModifierKey) {
//     if action != glfw.Press && action != glfw.Repeat {
//         return
//     }

//     switch key {
//     case glfw.KeyTab:
//         l.focusField = (l.focusField + 1) % 2
//     case glfw.KeyBackspace:
//         if l.focusField == 0 && len(l.username) > 0 {
//             l.username = l.username[:len(l.username)-1]
//         } else if l.focusField == 1 && len(l.password) > 0 {
//             l.password = l.password[:len(l.password)-1]
//         }
//     case glfw.KeyEnter:
//         l.tryLogin()
//     case glfw.KeyEscape:
//         l.done = true
//     }
// }

// func (l *LoginWindow) handleChar(w *glfw.Window, char rune) {
//     if l.focusField == 0 {
//         l.username += string(char)
//     } else {
//         l.password += string(char)
//     }
// }

// func (l *LoginWindow) tryLogin() {
//     if l.username == "" || l.password == "" {
//         l.errMsg = "用户名和密码不能为空"
//         return
//     }

//     resp, err := grpcClient.Login(l.username, l.password)
//     if err != nil {
//         l.errMsg = fmt.Sprintf("登录失败: %v", err)
//         return
//     }

//     if resp.Token == "" {
//         l.errMsg = "登录失败: 服务器未返回有效令牌"
//         return
//     }

//     // Store token and user info
//     grpcClient.SetToken(resp.Token)
//     // Could store user info here if needed
//     // user := resp.User

//     l.success = true
//     l.done = true
// }

// func (l *LoginWindow) Render() {
//     l.win.MakeContextCurrent()
//     gl.ClearColor(0.2, 0.3, 0.3, 1.0)
//     gl.Clear(gl.COLOR_BUFFER_BIT)

//     // TODO: Render login form using OpenGL
//     // For now, just display in the window title
//     status := fmt.Sprintf("Login [%s] %s",
//         map[int]string{0: "Username", 1: "Password"}[l.focusField],
//         l.errMsg)
//     l.win.SetTitle(status)

//     l.win.SwapBuffers()
// }

// func (l *LoginWindow) ShouldClose() bool {
//     return l.win.ShouldClose() || l.done
// }

// func (l *LoginWindow) Close() {
//     l.win.Destroy()
// }

// func ShowLogin() bool {
//     login, err := NewLoginWindow()
//     if err != nil {
//         log.Printf("Failed to create login window: %v", err)
//         return false
//     }
//     defer login.Close()

//     for !login.ShouldClose() {
//         login.Render()
//         glfw.PollEvents()
//     }

//     return login.success
// }
