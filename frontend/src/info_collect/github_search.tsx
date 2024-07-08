import React, { useState, useEffect } from "react";
import { Input, Row, Col, Table, Typography, Button, Modal } from "antd";
import type { ColumnsType, TablePaginationConfig } from "antd/es/table";
import { useNavigate } from "react-router-dom";

import { error } from "../utils/Message";
import { useSessionStorage } from "../utils/Cache";
import { GithubSearchPoc, GithubSearchCode } from "../../wailsjs/go/main/App";
import { BrowserOpenURL } from "../../wailsjs/runtime/runtime";
import CodeHighlightBox from "../utils/Code";

const { Link } = Typography;

const GithubSearchPage: React.FC = () => {
    const [data, setData, removeData] = useSessionStorage("githubSearchData", "");
    const [searchTerm, setSearchTerm] = useSessionStorage("githubSearchQuery", "");
    const [currentPagination, setCurrentPagination] = useSessionStorage("githubSearchPagetion", { current: 1 });
    const [isLoading, setIsLoading] = useState(false);
    const [currentId, setCurrentId] = useState<number>(1);

    const [isModalVisible, setIsModalVisible] = useState(false); // Code 模态框
    const [githubCode, setgithubCode] = useSessionStorage("githubCode", "");
    const [codeFileExtension, setCodeFileExtension] = useSessionStorage("codeFileExtension", ""); // 用于存储代码文件的后缀
    const [navigateState, setNavigateState] = useState(false); // 用于控制跳转
    const [nucleiNav, setNucleiNav] = useSessionStorage("nucleiNav", ""); // 用于控制跳转
    const [nucleiPlugin, setNucleiPlugin] = useSessionStorage("nucleiInput", "");

    interface DataType {
        full_name: string;
        html_url: string;
        pushed_at: string;
        file_name: string;
    }

    const handleSearch = async (query: string) => {
        if (query === "") {
            error("请输入正确查询语句");
            return;
        }

        removeData("githubSearchData");
        setCurrentId(1);
        setIsLoading(true);
        setSearchTerm(query);

        // 调用GithubSearchPoc接口
        GithubSearchPoc(query).then((response: any) => {
            console.log(response);
            if (response.code === 500) {
                error(response.msg);
                setIsLoading(false);
            } else {
                // 如果 msg 为空，则返回
                if (response.msg.length === 0) {
                    error("No data found");
                    setIsLoading(false);
                    return;
                }

                const formatDataForTable = response.msg.map((item: any, index: any) => ({
                    full_name: item.full_name,
                    html_url: item.html_url,
                    pushed_at: item.pushed_at,
                    file_name: item.file_name,
                }));
                setData(formatDataForTable);
            }
            setIsLoading(false);
        });
    };

    // 获取 code
    const handleGithubSearchCode = (html_url: string) => {
        if (html_url) {
            GithubSearchCode(html_url).then((response: any) => {
                if (response.code === 200) {
                    const newStr = response.msg.replace(/\n/g, "");
                    const binaryString = window.atob(newStr);
                    const bytes = new Uint8Array(newStr.length);
                    for (let i = 0; i < binaryString.length; i++) {
                        bytes[i] = binaryString.charCodeAt(i);
                    }

                    // 创建解码器并解码二进制数据
                    const decoder = new TextDecoder("utf-8");
                    const decodedString = decoder.decode(bytes);
                    console.log(decodedString);
                    setgithubCode(decodedString);

                    // 获取html_url中的文件名后缀
                    const url = new URL(html_url);
                    const pathname = url.pathname;
                    const filename = pathname.split("/").pop(); // 获取url路径中的最后一部分作为文件名
                    const extension = filename?.split(".").pop(); // 通过点号分隔文件名，取最后一部分作为后缀
                    setCodeFileExtension(extension);
                } else {
                    error("Error fetching Code: " + response.msg);
                }
            });
        } else {
            error("Path not found for Code: " + html_url);
        }
    };

    const showModalGetCode = (html_url: string) => {
        handleGithubSearchCode(html_url);
        setIsModalVisible(true);
    };

    const handleOk = () => {
        setIsModalVisible(false);
    };

    const handleCancel = () => {
        setIsModalVisible(false);
    };

    // 生成 nuclei 插件 --------------------------------------------------------------------------------------
    const handleNavigateToWebPluginsGen = () => {
        setIsModalVisible(false); // 关闭模态框
        setNavigateState(true); // 设置 navigateState 为 true，用于控制跳转
        setNucleiNav("nuclei"); // 设置 nucleiNav 为 nuclei，用于控制跳转
        setNucleiPlugin(githubCode); // 将 nucleiTemplate 写入 nucleiInput
    };

    // Separate useEffect to detect changes to nucleiPlugin
    const navigate = useNavigate();
    useEffect(() => {
        if (githubCode != "" && nucleiPlugin === githubCode && navigateState) {
            navigate("/web_plugins_gen");
        }
    }, [nucleiPlugin, nucleiNav]);

    const columns: ColumnsType<DataType> = [
        { title: "ID", key: "id", render: (_text, record, index) => `${index + 1}` },
        { title: "Full Name", dataIndex: "full_name", key: "full_name" },
        {
            title: "HTML URL",
            dataIndex: "html_url",
            key: "html_url",
            render: (url) => <Link onClick={() => BrowserOpenURL(url)}>{url}</Link>,
        },
        {
            title: "File Name",
            dataIndex: "file_name",
            key: "file_name",
            render: (_text, record, index) => <Typography.Text onClick={() => showModalGetCode(record.html_url)}>{record.file_name}</Typography.Text>,
        },
        { title: "Last Push", dataIndex: "pushed_at", key: "pushed_at" },
    ];

    // 分页和筛选 ------------------------------------------------------------------------------------------
    // 当分页或筛选条件改变时更新缓存
    const handleTableChange = (pagination: TablePaginationConfig) => {
        setCurrentPagination(pagination);
    };

    useEffect(() => {
        handleTableChange(currentPagination);
    }, [currentPagination]);

    return (
        <div style={{ padding: "20px" }}>
            <h1 style={{ textAlign: "left" }}>Github Search</h1>
            <Row gutter={8} style={{ marginBottom: "20px" }}>
                <Col span={21}>
                    <Input
                        value={searchTerm} // 使用 searchTerm
                        onChange={(e) => setSearchTerm(e.target.value)} // 更新 searchTerm
                        onPressEnter={() => handleSearch(searchTerm)} // 使用 searchTerm 进行查询
                        placeholder="Search"
                    />
                </Col>
                <Col span={3}>
                    <Button type="primary" block onClick={() => handleSearch(searchTerm)} loading={isLoading}>
                        Query
                    </Button>
                </Col>
            </Row>

            <Table dataSource={data} columns={columns} rowKey="id" pagination={{ ...currentPagination, defaultPageSize: 100 }} onChange={handleTableChange} />

            <Modal
                title="Code"
                open={isModalVisible}
                onOk={handleOk}
                onCancel={handleCancel}
                width={800}
                centered
                footer={[
                    <Button key="navigate" type="primary" onClick={handleNavigateToWebPluginsGen}>
                        生成 Web 插件
                    </Button>,
                ]}>
                <CodeHighlightBox language={codeFileExtension} code={githubCode} />
            </Modal>
        </div>
    );
};

export default GithubSearchPage;
