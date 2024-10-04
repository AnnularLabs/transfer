// app/root.tsx
import { ChakraProvider } from "@chakra-ui/react";
import { Outlet } from "@remix-run/react";

export default function App() {
  return (
    <ChakraProvider>
      <Outlet />
    </ChakraProvider>
  );
}
