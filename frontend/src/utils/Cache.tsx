import { useState, useEffect } from "react";

function useLocalStorage(key: string, initialValue: string) {
    // 获取存储值，或者如果没有找到则返回 initialValue
    const storedValue = localStorage.getItem(key);
    const initial = storedValue ? JSON.parse(storedValue) : initialValue;

    // 设置状态
    const [value, setValue] = useState(initial);

    // 当状态更改时，自动更新 localStorage
    useEffect(() => {
        localStorage.setItem(key, JSON.stringify(value));
    }, [key, value]);

    return [value, setValue];
}

function useSessionStorage(key: string, initialValue: any) {
    // Get stored value or use initialValue if not found
    const storedValue = sessionStorage.getItem(key);
    const initial = storedValue ? JSON.parse(storedValue) : initialValue;

    // Set state
    const [value, setValue] = useState(initial);

    // Automatically update sessionStorage when state changes
    useEffect(() => {
        sessionStorage.setItem(key, JSON.stringify(value));
    }, [key, value]);

    // Remove item from sessionStorage
    const remove = () => {
        try {
            sessionStorage.removeItem(key);
            setValue(initialValue); // Reset to initialValue after removal
        } catch (e) {
            console.error(`Error removing item with key "${key}" from sessionStorage`, e);
        }
    };

    // Clear all items from sessionStorage
    const clear = () => {
        try {
            sessionStorage.clear();
        } catch (e) {
            console.error("Error clearing sessionStorage", e);
        }
    };

    return [value, setValue, remove, clear];
}

export default useSessionStorage;
export { useLocalStorage, useSessionStorage };