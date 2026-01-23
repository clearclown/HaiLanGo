'use client';

import { AppLayout } from '@/components/layout';
import { Button, Card, CardContent, CardHeader } from '@/components/ui';
import { apiClient } from '@/lib/api/client';
import { useCallback, useEffect, useRef, useState } from 'react';

interface ConversationSession {
  sessionId: string;
  websocketUrl: string;
  targetLanguage: string;
  nativeLanguage: string;
}

interface Message {
  id: string;
  type: 'user' | 'assistant';
  text: string;
  timestamp: Date;
}

export default function ConversationPage() {
  const [isConnected, setIsConnected] = useState(false);
  const [isRecording, setIsRecording] = useState(false);
  const [session, setSession] = useState<ConversationSession | null>(null);
  const [messages, setMessages] = useState<Message[]>([]);
  const [error, setError] = useState<string | null>(null);
  const [targetLanguage, setTargetLanguage] = useState('ru');
  const [nativeLanguage, setNativeLanguage] = useState('ja');
  const [audioLevel, setAudioLevel] = useState(0);

  const wsRef = useRef<WebSocket | null>(null);
  const mediaRecorderRef = useRef<MediaRecorder | null>(null);
  const audioContextRef = useRef<AudioContext | null>(null);
  const analyserRef = useRef<AnalyserNode | null>(null);
  const messagesEndRef = useRef<HTMLDivElement>(null);

  // Scroll to bottom when messages change
  useEffect(() => {
    messagesEndRef.current?.scrollIntoView({ behavior: 'smooth' });
  }, [messages]);

  // Start conversation session
  const startSession = useCallback(async () => {
    try {
      setError(null);
      const response = await apiClient.conversation.start(targetLanguage, nativeLanguage);
      setSession({
        sessionId: response.session_id,
        websocketUrl: response.websocket_url,
        targetLanguage,
        nativeLanguage,
      });

      // Connect WebSocket
      const wsUrl = `${window.location.protocol === 'https:' ? 'wss:' : 'ws:'}//${window.location.host}${response.websocket_url}`;
      const ws = new WebSocket(wsUrl);

      ws.onopen = () => {
        setIsConnected(true);
        addMessage('assistant', 'Session started. Click the microphone to speak.');
      };

      ws.onmessage = (event) => {
        const data = JSON.parse(event.data);
        handleServerMessage(data);
      };

      ws.onclose = () => {
        setIsConnected(false);
        setSession(null);
      };

      ws.onerror = () => {
        setError('WebSocket connection failed');
        setIsConnected(false);
      };

      wsRef.current = ws;
    } catch (err) {
      console.error('Failed to start session:', err);
      setError('Failed to start conversation session');
    }
  }, [targetLanguage, nativeLanguage]);

  // Stop conversation session
  const stopSession = useCallback(async () => {
    if (session) {
      try {
        await apiClient.conversation.stop(session.sessionId);
      } catch (err) {
        console.error('Failed to stop session:', err);
      }
    }

    if (wsRef.current) {
      wsRef.current.close();
      wsRef.current = null;
    }

    setIsConnected(false);
    setSession(null);
    setIsRecording(false);
    stopRecording();
  }, [session]);

  // Handle server messages
  const handleServerMessage = useCallback((data: { type: string; text?: string; audio?: string; error?: string }) => {
    switch (data.type) {
      case 'response.text.delta':
      case 'response.audio_transcript.delta':
        if (data.text) {
          // Append to last assistant message or create new one
          setMessages(prev => {
            const last = prev[prev.length - 1];
            if (last && last.type === 'assistant') {
              return [...prev.slice(0, -1), { ...last, text: last.text + data.text }];
            }
            return [...prev, { id: Date.now().toString(), type: 'assistant', text: data.text || '', timestamp: new Date() }];
          });
        }
        break;

      case 'response.audio.delta':
        if (data.audio) {
          playAudio(data.audio);
        }
        break;

      case 'response.done':
        // Response complete
        break;

      case 'error':
        setError(data.error || 'An error occurred');
        break;

      case 'session.created':
        // Session created, ready to use
        break;
    }
  }, []);

  // Add message to chat
  const addMessage = (type: 'user' | 'assistant', text: string) => {
    setMessages(prev => [...prev, {
      id: Date.now().toString(),
      type,
      text,
      timestamp: new Date(),
    }]);
  };

  // Start recording
  const startRecording = async () => {
    try {
      const stream = await navigator.mediaDevices.getUserMedia({ audio: true });

      // Setup audio context for visualization
      audioContextRef.current = new AudioContext();
      analyserRef.current = audioContextRef.current.createAnalyser();
      const source = audioContextRef.current.createMediaStreamSource(stream);
      source.connect(analyserRef.current);
      analyserRef.current.fftSize = 256;

      // Start visualization loop
      const updateLevel = () => {
        if (analyserRef.current && isRecording) {
          const dataArray = new Uint8Array(analyserRef.current.frequencyBinCount);
          analyserRef.current.getByteFrequencyData(dataArray);
          const average = dataArray.reduce((a, b) => a + b) / dataArray.length;
          setAudioLevel(average / 255);
          requestAnimationFrame(updateLevel);
        }
      };
      updateLevel();

      // Setup MediaRecorder
      const mediaRecorder = new MediaRecorder(stream, {
        mimeType: 'audio/webm;codecs=opus',
      });

      mediaRecorder.ondataavailable = async (event) => {
        if (event.data.size > 0 && wsRef.current?.readyState === WebSocket.OPEN) {
          // Convert to base64 and send
          const reader = new FileReader();
          reader.onload = () => {
            const base64 = (reader.result as string).split(',')[1];
            wsRef.current?.send(JSON.stringify({
              type: 'audio',
              audio: base64,
            }));
          };
          reader.readAsDataURL(event.data);
        }
      };

      mediaRecorder.start(100); // Send chunks every 100ms
      mediaRecorderRef.current = mediaRecorder;
      setIsRecording(true);

    } catch (err) {
      console.error('Failed to start recording:', err);
      setError('Failed to access microphone');
    }
  };

  // Stop recording
  const stopRecording = () => {
    if (mediaRecorderRef.current) {
      mediaRecorderRef.current.stop();
      mediaRecorderRef.current.stream.getTracks().forEach(track => track.stop());
      mediaRecorderRef.current = null;
    }

    if (audioContextRef.current) {
      audioContextRef.current.close();
      audioContextRef.current = null;
    }

    setIsRecording(false);
    setAudioLevel(0);
  };

  // Toggle recording
  const toggleRecording = () => {
    if (isRecording) {
      stopRecording();
      // Add user message placeholder
      addMessage('user', '[Voice message]');
    } else {
      startRecording();
    }
  };

  // Play audio from base64
  const playAudio = async (base64Audio: string) => {
    try {
      if (!audioContextRef.current) {
        audioContextRef.current = new AudioContext();
      }

      const binaryString = atob(base64Audio);
      const bytes = new Uint8Array(binaryString.length);
      for (let i = 0; i < binaryString.length; i++) {
        bytes[i] = binaryString.charCodeAt(i);
      }

      const audioBuffer = await audioContextRef.current.decodeAudioData(bytes.buffer);
      const source = audioContextRef.current.createBufferSource();
      source.buffer = audioBuffer;
      source.connect(audioContextRef.current.destination);
      source.start();
    } catch (err) {
      console.error('Failed to play audio:', err);
    }
  };

  // Cleanup on unmount
  useEffect(() => {
    return () => {
      stopSession();
    };
  }, [stopSession]);

  const languages = [
    { code: 'ru', name: 'Russian', flag: '\u{1F1F7}\u{1F1FA}' },
    { code: 'zh', name: 'Chinese', flag: '\u{1F1E8}\u{1F1F3}' },
    { code: 'ja', name: 'Japanese', flag: '\u{1F1EF}\u{1F1F5}' },
    { code: 'en', name: 'English', flag: '\u{1F1EC}\u{1F1E7}' },
    { code: 'es', name: 'Spanish', flag: '\u{1F1EA}\u{1F1F8}' },
    { code: 'fr', name: 'French', flag: '\u{1F1EB}\u{1F1F7}' },
    { code: 'de', name: 'German', flag: '\u{1F1E9}\u{1F1EA}' },
    { code: 'fa', name: 'Persian', flag: '\u{1F1EE}\u{1F1F7}' },
    { code: 'he', name: 'Hebrew', flag: '\u{1F1EE}\u{1F1F1}' },
    { code: 'tr', name: 'Turkish', flag: '\u{1F1F9}\u{1F1F7}' },
  ];

  return (
    <AppLayout>
      <div className="container-app py-6 lg:py-8 max-w-4xl mx-auto">
        {/* Header */}
        <div className="mb-8">
          <h1 className="text-3xl font-bold text-gray-900 mb-2">AI Voice Conversation</h1>
          <p className="text-gray-600">Practice speaking with an AI language tutor</p>
        </div>

        {/* Language Selection */}
        {!isConnected && (
          <Card className="mb-6">
            <CardHeader>
              <h2 className="text-xl font-semibold">Language Settings</h2>
            </CardHeader>
            <CardContent className="space-y-4">
              <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                <div>
                  <label htmlFor="target-language" className="block text-sm font-medium text-gray-700 mb-2">
                    Learning Language
                  </label>
                  <select
                    id="target-language"
                    value={targetLanguage}
                    onChange={(e) => setTargetLanguage(e.target.value)}
                    className="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-blue-500"
                  >
                    {languages.map((lang) => (
                      <option key={lang.code} value={lang.code}>
                        {lang.flag} {lang.name}
                      </option>
                    ))}
                  </select>
                </div>

                <div>
                  <label htmlFor="native-language" className="block text-sm font-medium text-gray-700 mb-2">
                    Native Language
                  </label>
                  <select
                    id="native-language"
                    value={nativeLanguage}
                    onChange={(e) => setNativeLanguage(e.target.value)}
                    className="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-blue-500"
                  >
                    {languages.map((lang) => (
                      <option key={lang.code} value={lang.code}>
                        {lang.flag} {lang.name}
                      </option>
                    ))}
                  </select>
                </div>
              </div>

              <Button onClick={startSession} className="w-full">
                Start Conversation
              </Button>
            </CardContent>
          </Card>
        )}

        {/* Error Display */}
        {error && (
          <div className="mb-4 p-4 bg-red-50 border border-red-200 rounded-lg text-red-700">
            {error}
            <button
              type="button"
              onClick={() => setError(null)}
              className="ml-2 text-red-500 hover:text-red-700"
            >
              Dismiss
            </button>
          </div>
        )}

        {/* Chat Area */}
        {isConnected && (
          <Card className="mb-6">
            <CardHeader className="flex flex-row items-center justify-between">
              <div>
                <h2 className="text-xl font-semibold">Conversation</h2>
                <p className="text-sm text-gray-500">
                  {languages.find(l => l.code === targetLanguage)?.flag} Learning {languages.find(l => l.code === targetLanguage)?.name}
                </p>
              </div>
              <Button variant="danger" size="sm" onClick={stopSession}>
                End Session
              </Button>
            </CardHeader>
            <CardContent>
              {/* Messages */}
              <div className="h-96 overflow-y-auto mb-4 p-4 bg-gray-50 rounded-lg space-y-4">
                {messages.length === 0 ? (
                  <div className="text-center text-gray-400 py-8">
                    Click the microphone button to start speaking
                  </div>
                ) : (
                  messages.map((msg) => (
                    <div
                      key={msg.id}
                      className={`flex ${msg.type === 'user' ? 'justify-end' : 'justify-start'}`}
                    >
                      <div
                        className={`max-w-[80%] px-4 py-2 rounded-lg ${
                          msg.type === 'user'
                            ? 'bg-blue-600 text-white'
                            : 'bg-white border border-gray-200 text-gray-800'
                        }`}
                      >
                        <p className="whitespace-pre-wrap">{msg.text}</p>
                        <p className={`text-xs mt-1 ${msg.type === 'user' ? 'text-blue-200' : 'text-gray-400'}`}>
                          {msg.timestamp.toLocaleTimeString()}
                        </p>
                      </div>
                    </div>
                  ))
                )}
                <div ref={messagesEndRef} />
              </div>

              {/* Recording Controls */}
              <div className="flex flex-col items-center space-y-4">
                {/* Audio Level Visualization */}
                {isRecording && (
                  <div className="w-full h-2 bg-gray-200 rounded-full overflow-hidden">
                    <div
                      className="h-full bg-blue-600 transition-all duration-100"
                      style={{ width: `${audioLevel * 100}%` }}
                    />
                  </div>
                )}

                {/* Microphone Button */}
                <button
                  type="button"
                  onClick={toggleRecording}
                  className={`w-20 h-20 rounded-full flex items-center justify-center transition-all ${
                    isRecording
                      ? 'bg-red-600 hover:bg-red-700 animate-pulse'
                      : 'bg-blue-600 hover:bg-blue-700'
                  }`}
                >
                  {isRecording ? (
                    <svg className="w-10 h-10 text-white" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                      <rect x="6" y="6" width="12" height="12" rx="2" fill="currentColor" />
                    </svg>
                  ) : (
                    <svg className="w-10 h-10 text-white" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                      <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M19 11a7 7 0 01-7 7m0 0a7 7 0 01-7-7m7 7v4m0 0H8m4 0h4m-4-8a3 3 0 01-3-3V5a3 3 0 116 0v6a3 3 0 01-3 3z" />
                    </svg>
                  )}
                </button>

                <p className="text-sm text-gray-500">
                  {isRecording ? 'Recording... Click to stop' : 'Click to speak'}
                </p>
              </div>
            </CardContent>
          </Card>
        )}

        {/* Tips */}
        <Card>
          <CardHeader>
            <h2 className="text-lg font-semibold">Tips for Better Practice</h2>
          </CardHeader>
          <CardContent>
            <ul className="space-y-2 text-gray-600">
              <li className="flex items-start">
                <span className="text-blue-600 mr-2">1.</span>
                Speak clearly and at a natural pace
              </li>
              <li className="flex items-start">
                <span className="text-blue-600 mr-2">2.</span>
                Use headphones for better audio quality
              </li>
              <li className="flex items-start">
                <span className="text-blue-600 mr-2">3.</span>
                Try to respond in the target language as much as possible
              </li>
              <li className="flex items-start">
                <span className="text-blue-600 mr-2">4.</span>
                Ask the AI for pronunciation help when needed
              </li>
            </ul>
          </CardContent>
        </Card>
      </div>
    </AppLayout>
  );
}
