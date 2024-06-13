import {
  Box,
  Button,
  Cell,
  FormField,
  Input,
  InputArea,
  Layout,
  Text,
  WixDesignSystemProvider,
} from "@wix/design-system";
import { embeddedScripts } from "@wix/app-market";
import "@wix/design-system/styles.global.css";
import React, { useEffect, useRef, useState } from "react";
import { withDashboard, useDashboard } from "@wix/dashboard-react";
import { getAppInstance } from "../../auth";
import { useWixModules } from "@wix/sdk-react";

export interface Message {
  role: "user" | "assistant";
  content: string;
}

export interface StreamData {
  content?: string;
}

interface UserMessage {
  chat_id?: string;
  message: string;
}

const API_URL = "http://localhost:8001";

const fetchExistingChat = async (instanceId: string) => {
  const response = await fetch(`${API_URL}/placement-chat`, {
    method: "GET",
    headers: {
      "Content-Type": "application/json",
      accept: "application/json",
      Authorization: getAppInstance(),
    },
  });

  if (response.ok) {
    const chatHistory: Message[] = await response.json();
    return chatHistory;
  }
};

function ChatApp() {
  const chatBoxRef = useRef<HTMLDivElement>(null);
  const { embedScript } = useWixModules(embeddedScripts);

  const scrollToBottom = () => {
    if (chatBoxRef.current) {
      chatBoxRef.current.scrollTop = chatBoxRef.current.scrollHeight;
    }
  };

  const { showToast } = useDashboard();
  const [chatHistory, setChatHistory] = useState<Message[]>([
    {
      role: "assistant",
      content: `Hi! I'm Diego, your assistant.
        I'm going to ask you a few questions to check your Spanish level. 
        Can you introduce yourself in Spanish?
        
        ¿Cómo te llamas? ¿De dónde eres? ¿Qué haces en tu tiempo libre?`,
    },
  ]);

  useEffect(() => {
    embedScript({});
    fetchExistingChat(getAppInstance())
      .then((chatHistory) => {
        if (!chatHistory?.length) {
          return;
        }
        setChatHistory(chatHistory);
        scrollToBottom();
      })
      .catch((err) => {
        console.error(err);
        showToast({
          message: "Failed to fetch chat history",
          type: "error",
        });
      });
  }, []);
  const [userMessage, setUserMessage] = useState("");

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setChatHistory((prevHistory) => [
      ...prevHistory,
      { role: "user", content: userMessage },
    ]);
    scrollToBottom();
    setUserMessage("");
    const userMessagePayload: UserMessage = {
      // chat_id: "123",
      message: userMessage,
    };

    const response = await fetch(`${API_URL}/placement-chat`, {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
        accept: "application/json",
        Authorization: getAppInstance(),
      },
      body: JSON.stringify(userMessagePayload),
    });

    const reader = response.body?.getReader();
    let assistantMessage: Message = { role: "assistant", content: "" };

    console.log({ assistantMessage });
    while (true && reader) {
      const { done, value } = await reader.read();
      if (done) break;

      const chunks = new TextDecoder().decode(value).split("\n");
      chunks.filter(Boolean).forEach((chunk) => {
        try {
          if (!chunk.startsWith("data:")) return;
          const dataChunk = chunk.slice(5);
          if (!dataChunk.startsWith("{") || !dataChunk.endsWith("}")) return;
          const jsonChunk: StreamData = JSON.parse(chunk.slice(5));
          if (jsonChunk.content) {
            assistantMessage.content = jsonChunk.content;
            setChatHistory((prevHistory) => {
              const lastMessage = prevHistory[prevHistory.length - 1];
              if (lastMessage.role === "assistant") {
                return prevHistory.slice(0, -1).concat(assistantMessage);
              } else {
                return prevHistory.concat(assistantMessage);
              }
            });
            scrollToBottom();
          }
        } catch (err) {
          console.error(err);
        }
      });
    }
  };

  return (
    <WixDesignSystemProvider>
      <Box
        align="center"
        verticalAlign="middle"
        height="100vh"
        backgroundColor="D10"
      >
        <Box
          direction="vertical"
          padding="20px"
          backgroundColor="white"
          borderRadius="6px"
          boxShadow="medium"
          width="600px"
        >
          <Text size="medium" weight="bold" marginBottom="20px">
            Chat with Diego
          </Text>
          <Box
            direction="vertical"
            overflow="auto"
            height="400px"
            marginBottom="20px"
            ref={chatBoxRef}
          >
            {chatHistory.map((message, index) => (
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
                alignSelf={
                  message.role === "assistant" ? "flex-start" : "flex-end"
                }
              >
                <Text color="white">
                  <strong>{message.role === "user" ? "You" : "Diego"}:</strong>{" "}
                  {message.content}
                </Text>
              </Box>
            ))}
          </Box>
          <form onSubmit={handleSubmit}>
            <Layout>
              <Cell span={10}>
                <FormField>
                  <InputArea
                    placeholder="Enter your message"
                    value={userMessage}
                    onChange={(e) => setUserMessage(e.target.value)}
                    autoGrow
                  />
                </FormField>
              </Cell>
              <Cell span={2}>
                <FormField>
                  <Button
                    type="submit"
                    skin="standard"
                    borderRadius="0 6px 6px 0"
                  >
                    Send
                  </Button>
                </FormField>
              </Cell>
            </Layout>
          </form>
        </Box>
      </Box>
    </WixDesignSystemProvider>
  );
}

export default withDashboard(ChatApp);
