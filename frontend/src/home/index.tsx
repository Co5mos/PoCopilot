import { Link } from "react-router-dom";

import ToolCard from "../utils/Card";

import "./HomePage.css";

function chunkArray(array: any[], chunkSize: number) {
    const results = [];
    while (array.length) {
        results.push(array.splice(0, chunkSize));
    }
    return results;
}

function HomePage() {
    const tools = [
        {
            title: "GithubAction",
            description: "使用 Github Action 发包",
            route: "/github_action",
        },
        {
            title: "GithubSearch",
            description: "Github检索",
            route: "/github_search",
        },
    ];

    const chunkedTools = chunkArray([...tools], 3);

    return (
        <div id="body">
            <h1>所有工具</h1>

            <table>
                <tbody>
                    {chunkedTools.map((toolsRow, rowIndex) => (
                        <tr key={rowIndex} className="list-tr">
                            {toolsRow.map((tool, index) => (
                                <td key={index} className="list-td">
                                    <Link to={tool.route}>
                                        <ToolCard title={tool.title} description={tool.description} />
                                    </Link>
                                </td>
                            ))}
                        </tr>
                    ))}
                </tbody>
            </table>
        </div>
    );
}

export default HomePage;
