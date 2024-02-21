import { Modal } from 'antd';

export const success = (msg: string) => {
    Modal.success({
        title: "Success",
        content: (
            <div>
                <p>{msg}</p>
            </div>
        ),
        onOk() {},
        centered: true,
    });
};

export const error = (msg: string) => {
    Modal.error({
        title: "Error",
        content: (
            <div>
                <p>{msg}</p>
            </div>
        ),
        onOk() {},
        centered: true,
    });
};
