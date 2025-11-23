"use client";

import type React from "react";

import { useState, useEffect } from "react";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import { Label } from "@/components/ui/label";
import { useToast } from "@/hooks/use-toast";
import { Toaster } from "@/components/ui/toaster";
import { ExternalLink, Send, AlertCircle } from "lucide-react";
import { Alert, AlertDescription, AlertTitle } from "@/components/ui/alert";

interface Comment {
  id: string;
  email: string;
  content: string;
  created_at: string;
}

interface User {
  id: string;
  email: string;
}

interface AuthResponse {
  token?: string;
  user?: User;
  id?: string;
  email?: string;
  created_at?: string;
}

const API_URL = process.env.NEXT_PUBLIC_API_URL || "";
console.log("API_URL s:", API_URL);

export default function GoAPIDemo() {
  const [comments, setComments] = useState<Comment[]>([]);
  const [isLoggedIn, setIsLoggedIn] = useState(false);
  const [token, setToken] = useState<string>("");
  const [user, setUser] = useState<User | null>(null);
  const [isLoading, setIsLoading] = useState(false);

  // Auth form states
  const [showLogin, setShowLogin] = useState(false);
  const [showSignup, setShowSignup] = useState(false);
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");

  // Comment form state
  const [commentContent, setCommentContent] = useState("");

  const { toast } = useToast();

  useEffect(() => {
    if (API_URL) {
      fetchComments();
    }
  }, []);

  const fetchComments = async () => {
    try {
      const response = await fetch(`${API_URL}/comments`);
      if (response.ok) {
        const data = await response.json();
        setComments(Array.isArray(data) ? data : []);
      }
    } catch (error) {
      console.error("[v0] Failed to fetch comments:", error);
    }
  };

  const handleLogin = async (e: React.FormEvent) => {
    e.preventDefault();
    setIsLoading(true);

    try {
      const response = await fetch(`${API_URL}/login`, {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ email, password }),
      });

      if (response.ok) {
        const data: AuthResponse = await response.json();
        if (data.token && data.user) {
          setToken(data.token);
          setUser(data.user);
          setIsLoggedIn(true);
          setShowLogin(false);
          toast({
            title: "Login successful",
            description: `Welcome back, ${data.user.email}!`,
          });
          setEmail("");
          setPassword("");
        }
      } else {
        toast({
          title: "Login failed",
          description: "Invalid credentials",
          variant: "destructive",
        });
      }
    } catch (error) {
      console.error("[v0] Login error:", error);
      toast({
        title: "Error",
        description: "Failed to login",
        variant: "destructive",
      });
    } finally {
      setIsLoading(false);
    }
  };

  const handleSignup = async (e: React.FormEvent) => {
    e.preventDefault();
    setIsLoading(true);

    try {
      const response = await fetch(`${API_URL}/create`, {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ email, password }),
      });

      if (response.ok) {
        const data: AuthResponse = await response.json();
        toast({
          title: "Account created",
          description: `Account created for ${data.email}. Please login.`,
        });
        setShowSignup(false);
        setShowLogin(true);
        setPassword("");
      } else {
        toast({
          title: "Signup failed",
          description: "Could not create account",
          variant: "destructive",
        });
      }
    } catch (error) {
      console.error("[v0] Signup error:", error);
      toast({
        title: "Error",
        description: "Failed to create account",
        variant: "destructive",
      });
    } finally {
      setIsLoading(false);
    }
  };

  const handlePostComment = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!commentContent.trim()) return;

    setIsLoading(true);

    try {
      const response = await fetch(`${API_URL}/comment`, {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
          Authorization: token,
        },
        body: JSON.stringify({ content: commentContent }),
      });

      if (response.ok) {
        toast({
          title: "Comment posted",
          description: "Your comment has been added successfully",
        });
        setCommentContent("");
        fetchComments();
      } else {
        toast({
          title: "Failed to post comment",
          description: "Could not add your comment",
          variant: "destructive",
        });
      }
    } catch (error) {
      console.error("[v0] Comment post error:", error);
      toast({
        title: "Error",
        description: "Failed to post comment",
        variant: "destructive",
      });
    } finally {
      setIsLoading(false);
    }
  };

  return (
    <div className="min-h-screen bg-background">
      <div className="container max-w-4xl mx-auto px-4 py-12">
        {/* Header */}
        <div className="text-center mb-12">
          <h1 className="text-5xl font-bold text-balance mb-8">Go API Demo</h1>

          {/* Navigation Links */}
          <div className="flex items-center justify-center gap-4 flex-wrap">
            <a
              href={`${API_URL}/resume`}
              target="_blank"
              rel="noopener noreferrer"
              className="inline-flex items-center gap-2 px-6 py-2 rounded-lg border border-border hover:bg-accent transition-colors"
            >
              Resume
              <ExternalLink className="w-4 h-4" />
            </a>
            <a
              href={`${API_URL}/swagger`}
              target="_blank"
              rel="noopener noreferrer"
              className="inline-flex items-center gap-2 px-6 py-2 rounded-lg border border-border hover:bg-accent transition-colors"
            >
              Swagger
              <ExternalLink className="w-4 h-4" />
            </a>
            <a
              href="https://github.com/HWCronicus"
              target="_blank"
              rel="noopener noreferrer"
              className="inline-flex items-center gap-2 px-6 py-2 rounded-lg border border-border hover:bg-accent transition-colors"
            >
              GitHub
              <ExternalLink className="w-4 h-4" />
            </a>
          </div>
        </div>

        {!API_URL && (
          <Alert className="mb-8" variant="destructive">
            <AlertCircle className="h-4 w-4" />
            <AlertTitle>API URL Not Configured</AlertTitle>
            <AlertDescription>
              Please set the{" "}
              <code className="font-mono font-semibold">
                NEXT_PUBLIC_API_URL
              </code>{" "}
              environment variable in the Vars section of the sidebar to connect
              to your Go API.
            </AlertDescription>
          </Alert>
        )}

        {/* Comments Section */}
        <Card className="mb-8">
          <CardHeader>
            <CardTitle>Comments</CardTitle>
            <CardDescription>Recent comments from the API</CardDescription>
          </CardHeader>
          <CardContent>
            <div className="space-y-4">
              {comments.length === 0 ? (
                <p className="text-muted-foreground text-center py-8">
                  No comments yet
                </p>
              ) : (
                comments.map((comment) => (
                  <div
                    key={comment.id}
                    className="p-4 rounded-lg border border-border bg-card/50"
                  >
                    <div className="flex items-start justify-between gap-4 mb-2">
                      <span className="text-sm font-medium text-primary">
                        {comment.email}
                      </span>
                      <span className="text-xs text-muted-foreground">
                        {new Date(comment.created_at).toLocaleDateString()}
                      </span>
                    </div>
                    <p className="text-foreground">{comment.content}</p>
                  </div>
                ))
              )}
            </div>
          </CardContent>
        </Card>

        {/* Auth Section */}
        {!isLoggedIn ? (
          <div className="flex gap-4 justify-center mb-8">
            <Button
              onClick={() => {
                setShowLogin(true);
                setShowSignup(false);
              }}
              variant="default"
              size="lg"
              disabled={!API_URL}
            >
              Log In
            </Button>
            <Button
              onClick={() => {
                setShowSignup(true);
                setShowLogin(false);
              }}
              variant="outline"
              size="lg"
              disabled={!API_URL}
            >
              Sign Up
            </Button>
          </div>
        ) : (
          <div className="text-center mb-8">
            <p className="text-lg text-muted-foreground mb-4">
              Logged in as{" "}
              <span className="text-primary font-medium">{user?.email}</span>
            </p>
            <Button
              variant="outline"
              onClick={() => {
                setIsLoggedIn(false);
                setToken("");
                setUser(null);
                toast({ title: "Logged out successfully" });
              }}
            >
              Log Out
            </Button>
          </div>
        )}

        {/* Login Form */}
        {showLogin && !isLoggedIn && (
          <Card className="mb-8">
            <CardHeader>
              <CardTitle>Log In</CardTitle>
              <CardDescription>
                Enter your credentials to log in
              </CardDescription>
            </CardHeader>
            <CardContent>
              <form onSubmit={handleLogin} className="space-y-4">
                <div className="space-y-2">
                  <Label htmlFor="login-email">Email</Label>
                  <Input
                    id="login-email"
                    type="email"
                    value={email}
                    onChange={(e) => setEmail(e.target.value)}
                    placeholder="admin@example.com"
                    required
                  />
                </div>
                <div className="space-y-2">
                  <Label htmlFor="login-password">Password</Label>
                  <Input
                    id="login-password"
                    type="password"
                    value={password}
                    onChange={(e) => setPassword(e.target.value)}
                    placeholder="••••••••"
                    required
                  />
                </div>
                <div className="flex gap-2">
                  <Button type="submit" disabled={isLoading}>
                    {isLoading ? "Logging in..." : "Log In"}
                  </Button>
                  <Button
                    type="button"
                    variant="ghost"
                    onClick={() => setShowLogin(false)}
                  >
                    Cancel
                  </Button>
                </div>
              </form>
            </CardContent>
          </Card>
        )}

        {/* Signup Form */}
        {showSignup && !isLoggedIn && (
          <Card className="mb-8">
            <CardHeader>
              <CardTitle>Sign Up</CardTitle>
              <CardDescription>Create a new account</CardDescription>
            </CardHeader>
            <CardContent>
              <form onSubmit={handleSignup} className="space-y-4">
                <div className="space-y-2">
                  <Label htmlFor="signup-email">Email</Label>
                  <Input
                    id="signup-email"
                    type="email"
                    value={email}
                    onChange={(e) => setEmail(e.target.value)}
                    placeholder="admin@example.com"
                    required
                  />
                </div>
                <div className="space-y-2">
                  <Label htmlFor="signup-password">Password</Label>
                  <Input
                    id="signup-password"
                    type="password"
                    value={password}
                    onChange={(e) => setPassword(e.target.value)}
                    placeholder="••••••••"
                    required
                  />
                </div>
                <div className="flex gap-2">
                  <Button type="submit" disabled={isLoading}>
                    {isLoading ? "Creating account..." : "Sign Up"}
                  </Button>
                  <Button
                    type="button"
                    variant="ghost"
                    onClick={() => setShowSignup(false)}
                  >
                    Cancel
                  </Button>
                </div>
              </form>
            </CardContent>
          </Card>
        )}

        {/* Comment Form (only when logged in) */}
        {isLoggedIn && (
          <Card>
            <CardHeader>
              <CardTitle>Post a Comment</CardTitle>
              <CardDescription>Share your thoughts</CardDescription>
            </CardHeader>
            <CardContent>
              <form onSubmit={handlePostComment} className="space-y-4">
                <div className="space-y-2">
                  <Label htmlFor="comment-content">Your Comment</Label>
                  <Input
                    id="comment-content"
                    value={commentContent}
                    onChange={(e) => setCommentContent(e.target.value)}
                    placeholder="Write something..."
                    required
                  />
                </div>
                <Button type="submit" disabled={isLoading}>
                  <Send className="w-4 h-4 mr-2" />
                  {isLoading ? "Posting..." : "Post Comment"}
                </Button>
              </form>
            </CardContent>
          </Card>
        )}
      </div>

      <Toaster />
    </div>
  );
}

