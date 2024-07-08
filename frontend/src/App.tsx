import React, { useState, useEffect, useCallback } from "react";
import { Layout, Menu } from "antd";
import { ToolOutlined, FileOutlined, SecurityScanOutlined, GithubOutlined, CompassOutlined, FileTextOutlined } from "@ant-design/icons";
import { SettingOutlined } from "@ant-design/icons";
import { HashRouter, Route, Routes, Link, useLocation } from "react-router-dom";

import "./App.css";

import HomePage from "./home/index";
import GithubActionPage from "./packet_sender/github_action/index";
import GithubSearchPage from "./info_collect/github_search";
import SettingsPage from "./settings/index";

import { type CSSObject, Global } from "@emotion/react";
import { createAppTheme, createAppStylesBaseline } from "@arwes/react";

const theme = createAppTheme();
const stylesBaseline = createAppStylesBaseline(theme);

const { Sider } = Layout;
const { SubMenu } = Menu;

type NavigationProps = {
    openKeys: string[];
    setOpenKeys: React.Dispatch<React.SetStateAction<string[]>>;
};

function getOpenKeysFromPathname(pathname: string): string[] {
    // 使用一个对象来映射 pathname 到 openKeys
    const pathnameToOpenKeysMap: { [key: string]: string[] } = {
        "/github_action": ["/packet_sender"],
        "/github_search": ["/info_collect"],
    };

    // 从映射中获取 openKeys，如果没有找到则返回一个空数组
    return pathnameToOpenKeysMap[pathname] || [];
}

function App() {
    const [collapsed, setCollapsed] = useState(false);
    const [openKeys, setOpenKeys] = useState<string[]>([]); // 在这里添加useState

    return (
        <HashRouter basename={"/"}>
            <Global styles={stylesBaseline as Record<string, CSSObject>} />
            <Layout className="main-content">
                <Sider
                    className="fixed-sider"
                    collapsible
                    collapsed={collapsed}
                    onCollapse={() => setCollapsed(!collapsed)}
                    style={{ position: "relative", display: "flex", flexDirection: "column" }}>
                    <Navigation openKeys={openKeys} setOpenKeys={setOpenKeys} />
                </Sider>

                <Layout className="main-content">
                    <Routes>
                        <Route path="/" element={<HomePage />}></Route>
                        <Route path="/github_action" element={<GithubActionPage />}></Route>
                        <Route path="/github_search" element={<GithubSearchPage />}></Route>
                        <Route path="/settings" element={<SettingsPage />}></Route>
                    </Routes>
                </Layout>
            </Layout>
        </HashRouter>
    );
}

function Navigation({ openKeys, setOpenKeys }: NavigationProps) {
    const location = useLocation();
    const MemoizedMenuItem = React.memo(Menu.Item);

    useEffect(() => {
        const newOpenKeys = getOpenKeysFromPathname(location.pathname);
        if (JSON.stringify(newOpenKeys) !== JSON.stringify(openKeys)) {
            setOpenKeys(newOpenKeys);
        }
    }, [location.pathname]);

    const handleOpenChange = useCallback((keys: React.Key[]) => {
        setOpenKeys(keys as string[]);
    }, []);

    return (
        <Menu
            theme="dark"
            defaultSelectedKeys={["/"]}
            mode="inline"
            selectedKeys={[location.pathname]}
            openKeys={openKeys}
            onOpenChange={handleOpenChange}
            style={{ flex: 1 }}>
            <MemoizedMenuItem key="/" style={{ border: "none", outline: "none" }} icon={<ToolOutlined />}>
                <Link to="/">所有工具</Link>
            </MemoizedMenuItem>

            <SubMenu key="/packet_sender" icon={<SecurityScanOutlined />} title="发包器">
                <MemoizedMenuItem key="/github_action" icon={<GithubOutlined />}>
                    <Link to="/github_action">Github Action</Link>
                </MemoizedMenuItem>
            </SubMenu>

            <SubMenu key="/info_collect" icon={<CompassOutlined />} title="信息收集">
                <MemoizedMenuItem key="/github_search">
                    <Link to="/github_search">GithubSearch</Link>
                </MemoizedMenuItem>
            </SubMenu>

            <MemoizedMenuItem key="/settings" icon={<SettingOutlined />}>
                <Link to="/settings">设置</Link>
            </MemoizedMenuItem>
        </Menu>
    );
}

export default App;
