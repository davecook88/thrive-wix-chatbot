import { Box, Text } from "@wix/design-system";
import { Message } from "../../dashboard/pages/page";
import { FC, useRef } from "react";
import React from "react";

interface ChatViewProps {
  messages: Message[];
}

export const ChatView: FC<ChatViewProps> = ({ messages }) => {
  const chatBoxRef = useRef<HTMLDivElement>(null);
  return (
    <Box
      direction="vertical"
      overflow="auto"
      height="400px"
      marginBottom="20px"
      background={"#f0f0f0"}
      borderRadius={"20px"}
      padding={"20px"}
      ref={chatBoxRef}
    >
      {messages
        .filter((m) => m.role !== "system")
        .map((message, index) => (
          <Box
            key={index}
            direction="horizontal"
            backgroundColor={
              message.role === "assistant" ? "	#39ff5a" : "	#218aff"
            }
            borderRadius="6px"
            padding="10px"
            marginBottom="10px"
            maxWidth="80%"
            alignSelf={message.role === "assistant" ? "flex-start" : "flex-end"}
          >
            <Text color="white">
              <strong>{message.role === "user" ? "You" : "Diego"}:</strong>{" "}
              {message.content}
            </Text>
          </Box>
        ))}
    </Box>
  );
};
