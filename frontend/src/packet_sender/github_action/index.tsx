import { Row, Col, Input, Button, Tabs } from "antd";
import React, { useState } from "react";

import "./../../App.css";
import { error } from "../../utils/Message";
import CodeHighlightBox from "../../utils/Code";
import useWebSocket from "../../utils/WebSocket";
import { useSessionStorage } from "../../utils/Cache";

import { SendGithubAction, GetGithubActionLog } from "../../../wailsjs/go/main/App";


// TODO 剔除目标URL最后的空行
// TODO 补全 GET 请求缺失的空行
// TODO 增加 Curl 的 redirect 参数
const GithubActionPage: React.FC = () => {
    const [inputOne, setInputOne] = useSessionStorage("githubActionInputOne", ""); // 使用 useSessionStorage
    const [inputTwo, setInputTwo] = useSessionStorage("githubActionInputTwo", ""); // 使用 useSessionStorage
    const [logs, setLogs, removeLogs] = useSessionStorage("actionLogs", "");
    const [activeTabKey, setActiveTabKey] = useState("1"); // Default to the first tab

    const handleTabChange = (key: string) => {
        setActiveTabKey(key);
    };

    const handleWebSocketMessage = (event: any) => {
        const newLog = JSON.parse(event.data);
        // 如果 newLog 为空行，则不添加到日志中
        if (newLog === "") {
            return;
        }
        setLogs((prevLogs: any) => `${prevLogs}\n${newLog}`);
    };

    // Using the useWebSocket hook
    const { initiateConnection } = useWebSocket(handleWebSocketMessage);

    const handleSend = async () => {
        if (!inputOne || !inputTwo) {
            error("需要输入 Raw Data 和 URL");
            return; // Stop the function from proceeding if either input is empty
        }

        setLogs(""); // 清除 logs
        removeLogs("actionLogs"); // 清除 githubActionInputOne
        initiateConnection(); // 连接 websocket

        const getProtocolAndHostFromUrl = (url: string): string => {
            const parsedUrl = new URL(url);
            return `${parsedUrl.protocol}//${parsedUrl.host}`;
        };

        const inputTwoArray: string[] = inputTwo.split("\n");
        let hostsArray: string[] = [];
        let erroneousUrls: string[] = [];

        hostsArray = inputTwoArray.map((url) => {
            try {
                return getProtocolAndHostFromUrl(url);
            } catch (err) {
                erroneousUrls.push(url); // Push the erroneous URL to the array
                return "";
            }
        });

        // If there are erroneous URLs, report them all at once and return
        if (erroneousUrls.length > 0) {
            error(`URL 解析失败: ${erroneousUrls.join(", ")}`);
            return;
        }

        const inputTwoArrayLength = hostsArray.length;

        SendGithubAction(inputOne, hostsArray).then((response) => {
            if (response.code === 500) {
                error(response.msg);
            } else {
                GetGithubActionLog(inputTwoArrayLength).then((logResponse) => {
                    if (logResponse.code === 500) {
                        error(logResponse.msg);
                    }
                });
            }
        });

        // Switch to the "日志" tab
        setActiveTabKey("2");
    };

    const tabsItems = [
        {
            label: "发包",
            key: "1",
            children: (
                <>
                    <Row gutter={12}>
                        <Col span={12}>
                            <Input.TextArea
                                value={inputOne}
                                onChange={(e) => setInputOne(e.target.value)}
                                placeholder={`GET / HTTP/1.1
Host: www.cip.cc
Upgrade-Insecure-Requests: 1
User-Agent: Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36
Accept: text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9
Accept-Encoding: gzip, deflate
Accept-Language: zh-CN,zh;q=0.9
Connection: close

`}
                                style={{ height: "70vh" }}
                            />
                        </Col>
                        <Col span={12}>
                            <Input.TextArea
                                value={inputTwo}
                                onChange={(e) => setInputTwo(e.target.value)}
                                placeholder="http://www.cip.cc
http://www.cip.cc"
                                style={{ height: "70vh" }}
                            />
                        </Col>
                    </Row>

                    <div style={{ marginTop: 16, textAlign: "center" }}>
                        <Button type="primary" onClick={handleSend}>
                            发送
                        </Button>
                    </div>
                </>
            ),
        },

        {
            label: "日志",
            key: "2",
            children: (
                <div className="text-left">
                    <CodeHighlightBox language="http" code={logs} />
                </div>
            ),
        },
    ];

    return (
        <div id="body">
            <h1 style={{ textAlign: "left" }}>GithubAction</h1>
            <Tabs type="card" items={tabsItems} activeKey={activeTabKey} onChange={handleTabChange} /> {/* Pass activeKey prop here */}
        </div>
    );
};

export default GithubActionPage;
