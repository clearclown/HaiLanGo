import React, { useState } from 'react';
import { useAudioPlayer } from '@/hooks/useAudioPlayer';

interface AudioPlayerProps {
  audioUrl: string;
}

/**
 * 音声プレイヤーコンポーネント
 */
export const AudioPlayer: React.FC<AudioPlayerProps> = ({ audioUrl }) => {
  const { playing, currentTime, duration, speed, setSpeed, togglePlayPause } =
    useAudioPlayer({ audioUrl });
  const [showSpeedMenu, setShowSpeedMenu] = useState(false);

  const speeds = [0.5, 0.75, 1.0, 1.25, 1.5, 2.0];

  const handleSpeedChange = (newSpeed: number) => {
    setSpeed(newSpeed);
    setShowSpeedMenu(false);
  };

  const formatTime = (seconds: number) => {
    const mins = Math.floor(seconds / 60);
    const secs = Math.floor(seconds % 60);
    return `${mins}:${secs.toString().padStart(2, '0')}`;
  };

  return (
    <div className="bg-white p-4 rounded-lg shadow">
      <div className="flex items-center justify-between mb-2">
        <div className="text-sm text-gray-600">
          {formatTime(currentTime)} / {formatTime(duration)}
        </div>
        <div className="relative">
          <button
            onClick={() => setShowSpeedMenu(!showSpeedMenu)}
            className="px-3 py-1 bg-gray-100 rounded hover:bg-gray-200"
            aria-label={`${speed}x`}
          >
            {speed}x
          </button>
          {showSpeedMenu && (
            <div className="absolute right-0 mt-2 bg-white border rounded shadow-lg z-10">
              {speeds.map((s) => (
                <button
                  key={s}
                  onClick={() => handleSpeedChange(s)}
                  className="block w-full px-4 py-2 text-left hover:bg-gray-100"
                  aria-label={`${s}x`}
                >
                  {s}x
                </button>
              ))}
            </div>
          )}
        </div>
      </div>

      <div className="flex items-center space-x-4">
        <button
          onClick={togglePlayPause}
          className="flex-shrink-0 w-12 h-12 flex items-center justify-center bg-blue-500 text-white rounded-full hover:bg-blue-600"
          aria-label={playing ? '一時停止' : '再生'}
        >
          {playing ? '⏸' : '▶'}
        </button>

        <div className="flex-1">
          <div className="h-2 bg-gray-200 rounded-full overflow-hidden">
            <div
              className="h-full bg-blue-500"
              style={{ width: `${(currentTime / duration) * 100}%` }}
            />
          </div>
        </div>

        <button
          className="flex-shrink-0 px-3 py-1 bg-gray-100 rounded hover:bg-gray-200"
          aria-label="繰り返し"
        >
          🔁
        </button>
      </div>
    </div>
  );
};
