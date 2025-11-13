import { useState, useRef, useCallback } from 'react';

interface UseAudioPlayerProps {
  audioUrl: string;
  onEnded?: () => void;
}

interface UseAudioPlayerReturn {
  playing: boolean;
  currentTime: number;
  duration: number;
  speed: number;
  play: () => void;
  pause: () => void;
  seek: (time: number) => void;
  setSpeed: (speed: number) => void;
  togglePlayPause: () => void;
}

/**
 * 音声プレイヤーのカスタムフック
 */
export function useAudioPlayer({
  audioUrl,
  onEnded,
}: UseAudioPlayerProps): UseAudioPlayerReturn {
  const audioRef = useRef<HTMLAudioElement | null>(null);
  const [playing, setPlaying] = useState(false);
  const [currentTime, setCurrentTime] = useState(0);
  const [duration, setDuration] = useState(0);
  const [speed, setSpeedState] = useState(1.0);

  // Audioオブジェクトの初期化
  const initAudio = useCallback(() => {
    if (!audioRef.current) {
      const audio = new Audio(audioUrl);

      audio.addEventListener('timeupdate', () => {
        setCurrentTime(audio.currentTime);
      });

      audio.addEventListener('loadedmetadata', () => {
        setDuration(audio.duration);
      });

      audio.addEventListener('ended', () => {
        setPlaying(false);
        if (onEnded) {
          onEnded();
        }
      });

      audio.playbackRate = speed;
      audioRef.current = audio;
    }
  }, [audioUrl, speed, onEnded]);

  const play = useCallback(() => {
    initAudio();
    if (audioRef.current) {
      audioRef.current.play();
      setPlaying(true);
    }
  }, [initAudio]);

  const pause = useCallback(() => {
    if (audioRef.current) {
      audioRef.current.pause();
      setPlaying(false);
    }
  }, []);

  const seek = useCallback((time: number) => {
    if (audioRef.current) {
      audioRef.current.currentTime = time;
      setCurrentTime(time);
    }
  }, []);

  const setSpeed = useCallback((newSpeed: number) => {
    setSpeedState(newSpeed);
    if (audioRef.current) {
      audioRef.current.playbackRate = newSpeed;
    }
  }, []);

  const togglePlayPause = useCallback(() => {
    if (playing) {
      pause();
    } else {
      play();
    }
  }, [playing, play, pause]);

  return {
    playing,
    currentTime,
    duration,
    speed,
    play,
    pause,
    seek,
    setSpeed,
    togglePlayPause,
  };
}
