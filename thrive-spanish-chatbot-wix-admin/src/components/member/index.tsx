import { Box, Text } from "@wix/design-system";
import { useWixModules } from "@wix/sdk-react";

import { FC, useEffect, useState } from "react";
import { members } from "@wix/members";
import { contacts } from "@wix/crm";
import React from "react";

export const DisplayMemberDetails: FC<{ member: members.Member }> = ({
  member,
}) => {
  const { getContact } = useWixModules(contacts);
  const [contact, setContact] = useState<contacts.Contact | null>(null);
  useEffect(() => {
    if (!member.contactId) {
      return;
    }
    getContact(member.contactId).then((contact) => {
      setContact(contact);
    });
  }, [member.contactId]);
  console.log("contact", contact);

  const getContactNotes = (contact: contacts.Contact) => {
    const notes = contact.info?.extendedFields?.items?.["custom.notes"]
      ?.replace(/\*/g, "\n")
      .split("\n")
      .filter(Boolean)
      .join("\n");

    return notes;
  };

  const getContactEvaluation = (contact: contacts.Contact) => {
    const evaluation = contact.info?.labelKeys?.items?.find((item) =>
      item.startsWith("custom.level-")
    );
    return evaluation?.replace("custom.level-", "").toUpperCase();
  };
  return (
    <Box
      padding="20px"
      borderRadius="20px"
      width="100%"
      height="max-content"
      backgroundColor="white"
    >
      <table>
        <tbody>
          <tr>
            <td>
              <Text weight="bold">Email</Text>
            </td>
            <td>
              <Text>{contact?.primaryInfo?.email || "No email available"}</Text>
            </td>
          </tr>
          <tr>
            <td>
              <Text weight="bold">Notes</Text>
            </td>
            <td>
              <Text>
                {contact ? getContactNotes(contact) : "No notes available"}
              </Text>
            </td>
          </tr>
          <tr>
            <td>
              <Text weight="bold">Evaluation</Text>
            </td>
            <td>
              <Text>
                {contact
                  ? getContactEvaluation(contact)
                  : "No evaluation available"}
              </Text>
            </td>
          </tr>
        </tbody>
      </table>
    </Box>
  );
};
