import { BleepsOnAnimator, Animated, FrameSVGCorners, Text, aa, aaVisibility, Animator } from "@arwes/react";

import { createAppTheme } from "@arwes/react";
const theme = createAppTheme();

type ToolCardProps = {
    title: string;
    description: string;
};

const ToolCard: React.FC<ToolCardProps> = ({ title, description }) => {
    return (
        <Animator merge combine manager="stagger">
            {/* Play the intro bleep when card appears. */}
            <BleepsOnAnimator transitions={{ entering: "intro" }} continuous />

            <Animated
                className="card"
                style={{
                    position: "relative",
                    display: "block",
                    maxWidth: "300px",
                    margin: theme.space([4, "auto"]),
                    padding: theme.space(8),
                    textAlign: "center",
                    width: 220,
                    height: 160,
                }}
                // Effects for entering and exiting animation transitions.
                animated={[aaVisibility(), aa("y", "2rem", 0)]}
                // Play bleep when the card is clicked.
            >
                {/* Frame decoration and shape colors defined by CSS. */}
                <style>{`
                    [data-name=bg] {
                        color: ${theme.colors.primary.deco(1)};
                    }
                    [data-name=line] {
                        color: ${theme.colors.primary.main(4)};
                    }
                `}</style>

                {/* .arwes-react-frames-framesvg  */}
                <Animator>
                    <FrameSVGCorners strokeWidth={2} />
                </Animator>

                <Animator>
                    <Text as="h1">{title}</Text>
                </Animator>

                <Animator>
                    <Text>{description}</Text>
                </Animator>
            </Animated>
        </Animator>
    );
};

export default ToolCard;
