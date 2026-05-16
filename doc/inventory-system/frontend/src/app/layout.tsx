import type { Metadata } from "next";
import "./globals.css";
import QueryProvider from "@/components/common/QueryProvider";
import { AuthProvider } from "@/contexts/AuthContext";
import { AntdRegistry } from "@ant-design/nextjs-registry";

export const metadata: Metadata = {
  title: "库存管理系统",
  description: "企业级库存管理系统",
};

export default function RootLayout({
  children,
}: Readonly<{
  children: React.ReactNode;
}>) {
  return (
    <html lang="zh-CN">
      <body className="antialiased bg-gray-50 min-h-screen">
        <AntdRegistry>
          <QueryProvider>
            <AuthProvider>
              {children}
            </AuthProvider>
          </QueryProvider>
        </AntdRegistry>
      </body>
    </html>
  );
}
