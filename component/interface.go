package component

/**
 *
 * @author: schbook
 * @email: seekerxu@163.com
 * @date: 2018/6/4
 *
 */

type Plugin interface {
	Register(p Plugin, priority int)
	OnConnect(clientId, userName, passWord string) bool
	OnMessage()
}
