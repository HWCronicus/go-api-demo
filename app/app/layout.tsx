import type React from "react";
import type { Metadata } from "next";
import { Poppins } from "next/font/google";
import "./globals.css";

const _poppins = Poppins({ weight: ["600"], subsets: ["latin"] });

export const metadata: Metadata = {
  title: "Go API Demo",
  description: "Interactive demo for a Go API with authentication and comments",
};

export default function RootLayout({
  children,
}: Readonly<{
  children: React.ReactNode;
}>) {
  return (
    <html lang="en">
      <body className={`font-sans antialiased`}>{children}</body>
    </html>
  );
}

