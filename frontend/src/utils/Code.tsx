import React, { useEffect, useRef } from "react";
import { Button, Row, Col, Tooltip } from "antd";
import { CopyToClipboard } from "react-copy-to-clipboard";
import { CopyOutlined } from "@ant-design/icons";

import hljs from "highlight.js";
import "highlight.js/styles/atom-one-dark.css"; // 选择您喜欢的样式

import "./utils.css";

const CodeHighlightBox: React.FC<{ language: string; code: string }> = ({ language, code }) => {
    const codeEl = useRef<HTMLElement | null>(null);

    useEffect(() => {
        if (codeEl.current) {
            // Unset the 'highlighted' marker if it exists
            if (codeEl.current.dataset.highlighted) {
                delete codeEl.current.dataset.highlighted;
            }

            // Apply highlight.js
            hljs.highlightElement(codeEl.current);

            // Set the 'highlighted' marker
            codeEl.current.dataset.highlighted = "true";
        }
    }, [code]);

    return (
        <div className="code-highlight-box">
            <Row align="middle">
                <Col className="language-label">
                    <span style={{ color: "hsla(200,84%,54.74%,1)" }}>{language}</span>
                </Col>
                <Col flex="auto"></Col>
                <Col>
                    <CopyToClipboard text={code}>
                        <Tooltip title="Copy code">
                            <Button type="text" icon={<CopyOutlined style={{ color: "hsla(200,84%,54.74%,1)" }} />} />
                        </Tooltip>
                    </CopyToClipboard>
                </Col>
            </Row>

            <div>
                <pre>
                    <code ref={codeEl} className={`language-${language}`}>
                        {code}
                    </code>
                </pre>
            </div>
        </div>
    );
};

export default CodeHighlightBox;
