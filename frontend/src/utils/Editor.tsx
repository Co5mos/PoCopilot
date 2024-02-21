import React, { useEffect, useState } from "react";
import { Button, Row, Col, Tooltip, Switch } from "antd";
import { CopyToClipboard } from "react-copy-to-clipboard";
import { CopyOutlined } from "@ant-design/icons";
import AceEditor from "react-ace";

import "ace-builds/src-noconflict/mode-json";
import "ace-builds/src-noconflict/mode-yaml";
import "ace-builds/src-noconflict/mode-plain_text"
import "ace-builds/src-noconflict/theme-one_dark";
import "ace-builds/src-noconflict/ext-language_tools"; // for autocompletion

const HighlightedEditableBox: React.FC<{
    language: string;
    code: any;
    placeholder: string;
    hight: string;
    onChange: (newCode: any) => void;
}> = ({ language, code, placeholder, hight, onChange }) => {
    const [editorHeight, setEditorHeight] = useState(hight);
    const [wrapEnabled, setWrapEnabled] = useState(true);

    useEffect(() => {
        if (Array.isArray(code)) {
            code = code.toString();
        }

        if (!code) {
            return;
        }

        // 计算代码的行数
        const numLines = code.split(/\r\n|\r|\n/).length;

        // 假设每一行的高度是16像素，并加上一些额外的高度（例如，编辑器的边框、内边距等）
        const newHeight = `${numLines * 18}px`;

        if (newHeight < hight) {
            setEditorHeight(hight);
        } else {
            setEditorHeight(newHeight);
        }
    }, [code]);

    return (
        <div>
            <Row align="middle">
                <Col className="language-label">
                    <span>{language}</span>
                </Col>
                <Col flex="auto" style={{ textAlign: "right", paddingRight: "1%" }}>
                    <Switch
                        checkedChildren="换行"
                        unCheckedChildren="不换行"
                        defaultChecked
                        onChange={(checked) => {
                            setWrapEnabled(checked);
                        }}
                    />
                </Col>
                <Col>
                    <CopyToClipboard text={code}>
                        <Tooltip title="Copy code">
                            <Button type="text" icon={<CopyOutlined style={{ color: "white" }} />} />
                        </Tooltip>
                    </CopyToClipboard>
                </Col>
            </Row>

            <div>
                <AceEditor
                    mode={language}
                    theme="one_dark"
                    value={typeof code === 'object' ? JSON.stringify(code, null, 2) : code}
                    onChange={onChange} // 传入一个函数，当编辑器内容改变时，会调用这个函数
                    name="code-editor" // 用于多个编辑器
                    editorProps={{ $blockScrolling: true }} // 防止报错
                    fontSize={14} // 字体大小
                    placeholder={placeholder}
                    showPrintMargin={false} // 不显示右侧的竖线
                    showGutter={true} // 显示行号
                    highlightActiveLine={true} // 当前行高亮
                    wrapEnabled={wrapEnabled} // 自动换行
                    setOptions={{
                        enableBasicAutocompletion: false, // 启用基本自动完成
                        enableLiveAutocompletion: false, // 启用实时自动补全
                        enableSnippets: true, // 启用代码段
                        showLineNumbers: true, // 显示行号
                        tabSize: 4,
                    }}
                    width="100%"
                    height={editorHeight} // 使用动态计算的高度
                />
            </div>
        </div>
    );
};

export default HighlightedEditableBox;
