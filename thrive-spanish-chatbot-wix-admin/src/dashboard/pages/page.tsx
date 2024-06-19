import {
  Box,
  WixDesignSystemProvider,
  Text,
  Divider,
} from "@wix/design-system";
import "@wix/design-system/styles.global.css";
import React, { useEffect, useMemo, useState } from "react";
import { withDashboard, useDashboard } from "@wix/dashboard-react";
import { getAppInstance } from "../../auth";
import { ChatView } from "../../components/chat";
import { members } from "@wix/members";

import { useWixModules } from "@wix/sdk-react";
import { DisplayMemberDetails } from "../../components/member";

export interface Message {
  role: "user" | "assistant" | "system";
  content: string;
  name?: string;
}

export interface SavedChatMessage {
  messages: Message[];
  memberId: string;
  memberName: string;
  lastUpdated: string;
}

export interface StreamData {
  content?: string;
}

interface UserMessage {
  chat_id?: string;
  message: string;
}

const API_URL = "https://thrive-chat-ba0bf.uc.r.appspot.com";

function ChatApp() {
  const { showToast } = useDashboard();
  const [chats, setChats] = useState<SavedChatMessage[]>([]);
  const [selectedChat, setSelectedChat] = useState<SavedChatMessage | null>(
    null
  );
  const [member, setMember] = useState<members.Member | null>(null);
  const { getMember } = useWixModules(members);

  useEffect(() => {
    fetch(`${API_URL}/admin/list-chats`, {
      headers: {
        Authorization: getAppInstance(),
      },
    })
      .then(
        (response) => response.json() as Promise<{ chats: SavedChatMessage[] }>
      )
      .then((data) => {
        setChats(data.chats);
      });
  }, []);

  const formatDate = (dateString: string) => {
    const date = new Date(dateString);
    return date.toLocaleString();
  };

  useEffect(() => {
    if (!selectedChat) {
      setMember(null);
      return;
    }
    getMember(selectedChat?.memberId).then((member) => {
      setMember(member);
    });
  }, [selectedChat?.memberId]);

  console.log("member", member);
  return (
    <WixDesignSystemProvider>
      <Box display="flex" height="100vh" backgroundColor="D10">
        <Box
          direction="vertical"
          width="25%"
          backgroundColor="D20"
          padding="20px"
          overflow="auto"
          boxShadow="1px 0 5px rgba(0,0,0,0.1)"
        >
          <Text size="small" weight="bold" marginBottom="20px" skin="primary">
            Chats
          </Text>
          <Box direction="vertical">
            {chats?.map((chat) => (
              <div onClick={() => setSelectedChat(chat)} key={chat.memberId}>
                <Box
                  key={chat.memberId}
                  padding="10px"
                  backgroundColor={selectedChat === chat ? "D80" : "D60"}
                  margin="10px 0"
                  borderRadius="8px"
                  cursor="pointer"
                  boxShadow="0 2px 4px rgba(0,0,0,0.1)"
                  direction="vertical"
                >
                  <Text size="medium" weight="bold">
                    {chat.memberName}
                  </Text>
                  <Text size="small" color="D40">
                    {formatDate(chat.lastUpdated)}
                  </Text>
                </Box>
              </div>
            ))}
          </Box>
        </Box>
        <Divider direction="vertical" />
        <Box direction="vertical" width="75%" padding="20px">
          <Box
            backgroundColor="D10"
            overflow="auto"
            direction="vertical"
            gap={2}
          >
            {selectedChat ? (
              <ChatView
                messages={selectedChat.messages?.filter((m) => !m.name)}
              />
            ) : (
              <Text size="small" color="D40">
                Select a chat to view messages
              </Text>
            )}
          </Box>
          <Box>{member && <DisplayMemberDetails member={member} />}</Box>{" "}
        </Box>
      </Box>
    </WixDesignSystemProvider>
  );
}

export default withDashboard(ChatApp);
