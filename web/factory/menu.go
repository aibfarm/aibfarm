package factory

// MakeMenu : make menu template string
func MakeMenu() []*MenuEntry {

	var menuList []*MenuEntry

	var researchToolSub []*SubMenu

	//!--研究工具
	//-子目录

	//--子目录
	researchToolSub = append(researchToolSub, &SubMenu{
		Href:  "/aibfarm.debug",
		Title: "调试窗口",
	})

	//主目录
	researchTool := &MenuEntry{
		Icon:    "fas fa-search-dollar",
		Style:   "color: darkblue",
		Name:    "研究工具",
		ID:      "research_tool",
		SubMenu: researchToolSub,
	}

	menuList = append(menuList, researchTool)

	return menuList
}
