import type { MetaFunction } from "@remix-run/node";
import { Flex, Box, Button, Input, Progress, Textarea, IconButton, Center, HStack, Text } from "@chakra-ui/react";
import { useState } from "react";

export const meta: MetaFunction = () => {
  return [
    { title: "New Remix App" },
    { name: "description", content: "Welcome to Remix!" },
  ];
};

export default function Index() {
  const [progress, setProgress] = useState(30); // 假设目前的进度为30%

  const handleTransfer = () => {
    // 模拟传输操作，可以在这里发起API请求并更新进度
    setProgress(70); // 示例：将进度更新为70%
  };

  return (
    <Flex width="100vw" height="100vh" direction={"column"} alignItems={"center"}>
      <Box width={"70%"}>
        <Text pt={15} fontSize="5xl" fontWeight="black" textColor="#319795" opacity={"80%"}>Transfer</Text>
      </Box>
      <Center width="100vw" height="100vh">
        <Box pl={20} pr={20} pt={3} borderRadius="xl" boxShadow="md" width="70%" height="85vh" mx="auto">
          <HStack spacing={4} mb={3}>
            <Input placeholder="PlayListId" mb={3} size="lg" textColor={"#A0AEC0"} borderColor="#319795" />
            <IconButton
              aria-label="search music"
              mb={3}
              size="lg"
              isRound
            />
          </HStack>
          <Box mb={3} p={4} borderWidth={1} height={"70%"} borderRadius="lg" borderColor="#319795">
            <p>Hello World</p>
          </Box>
          <HStack spacing={4} pt={3} mb={3}>
            <Button
              colorScheme="#ffffff"
              borderWidth={1}
              borderColor="#319795"
              textColor="#319795"
              onClick={handleTransfer}
              size="lg"
              height="48px" // 设置按钮的高度
            >
              Transfer Spotify
            </Button>
            <Box borderWidth={1} borderColor="#319795" borderRadius="lg" pl={5} pr={5} flex="1" height="48px" display="flex" alignItems="center" justifyContent="center">
              <Progress
                value={progress}
                size="lg"
                colorScheme="green"
                hasStripe
                borderRadius="md"
                height="30%" // 让进度条填满 Box 的 50% 高度
                width="100%" // 让进度条填满 Box 的宽度
              />
            </Box>
          </HStack>
        </Box>
      </Center>
    </Flex>


  );
}

