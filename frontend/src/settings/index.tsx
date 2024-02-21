import { useEffect, useState } from "react";
import { Input, Select, Tooltip, Button, Modal, Form, Switch } from "antd";
import { InfoCircleOutlined, EyeTwoTone, EyeInvisibleOutlined } from "@ant-design/icons";
import { ReadConfig, WriteConfig } from "../../wailsjs/go/main/App";

function SettingsPage() {
    const [config, setConfig] = useState<any | null>(null);

    // 读取配置数据
    useEffect(() => {
        ReadConfig().then(setConfig);
    }, []);

    // 写配置数据
    const [isModalOpen, setIsModalOpen] = useState(false);
    const [modalMessage, setModalMessage] = useState("");

    const handleSubmit = async () => {
        // 从输入框和其他组件中读取配置数据
        const configData = {
            GithubToken: config?.GithubToken,
            Owner: config?.Owner,
            RepoName: config?.RepoName,
            Proxy: config?.Proxy,
            ProxySwitch: config?.ProxySwitch,
        };

        // 调用 Go 方法并传递配置数据
        const result = await WriteConfig(configData);
        // 设置模态框的消息并显示它
        setModalMessage(result);
        setIsModalOpen(true);
    };

    return (
        <div id="body">
            <Form labelCol={{ span: 3 }} layout="horizontal">
                {/* Group 1 */}
                <h6 style={{ textAlign: "left" }}>Github 配置</h6>

                <Form.Item label={<span style={{ color: "white" }}>Github Token</span>} style={{ marginBottom: "12px" }}>
                    <Input.Password
                        placeholder="Github Token"
                        value={config?.GithubToken}
                        onChange={(e) => setConfig({ ...config, GithubToken: e.target.value })}
                        iconRender={(visible) => (visible ? <EyeTwoTone /> : <EyeInvisibleOutlined />)}
                        suffix={
                            <Tooltip title="Extra information">
                                <InfoCircleOutlined style={{ color: "rgba(0,0,0,.45)" }} />
                            </Tooltip>
                        }
                    />
                </Form.Item>

                <Form.Item label={<span style={{ color: "white" }}>仓库拥有者</span>} style={{ marginBottom: "12px" }}>
                    <Input
                        placeholder="仓库拥有者"
                        value={config?.Owner}
                        onChange={(e) => setConfig({ ...config, Owner: e.target.value })}
                    />
                </Form.Item>

                <Form.Item label={<span style={{ color: "white" }}>仓库名称</span>}>
                    <Input
                        placeholder="仓库名称"
                        value={config?.RepoName}
                        onChange={(e) => setConfig({ ...config, RepoName: e.target.value })}
                    />
                </Form.Item>

                {/* Group 6 */}
                <h6 style={{ textAlign: "left" }}>其他配置</h6>

                <Form.Item label={<span style={{ color: "white" }}>Proxy</span>} style={{ marginBottom: "12px" }}>
                    <Input
                        placeholder="Proxy, eg: http://127.0.0.1:8080"
                        value={config?.Proxy}
                        onChange={(e) => setConfig({ ...config, Proxy: e.target.value })}
                    />
                </Form.Item>

                <Form.Item label={<span style={{ color: "white" }}>Proxy 开关</span>} style={{ marginBottom: "12px" }}>
                    <Switch
                        style={{ float: "left" }}
                        checkedChildren="开启"
                        unCheckedChildren="关闭"
                        checked={config?.ProxySwitch}
                        onChange={(checked) => {
                            setConfig({ ...config, ProxySwitch: checked });
                        }}
                    />
                </Form.Item>

                <Modal title="提示" open={isModalOpen} onOk={() => setIsModalOpen(false)} onCancel={() => setIsModalOpen(false)} centered>
                    <p>{modalMessage}</p>
                </Modal>

                <Button type="primary" onClick={handleSubmit} style={{ marginTop: "20px" }}>
                    提交
                </Button>
            </Form>
        </div>
    );
}

export default SettingsPage;
