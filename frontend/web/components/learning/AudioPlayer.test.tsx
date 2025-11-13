import { describe, it, expect, vi, beforeEach } from 'vitest';
import { render, screen, fireEvent } from '@testing-library/react';
import { AudioPlayer } from './AudioPlayer';

// Audioのモック
const mockPlay = vi.fn();
const mockPause = vi.fn();

global.Audio = vi.fn().mockImplementation(() => ({
  play: mockPlay,
  pause: mockPause,
  currentTime: 0,
  duration: 100,
  playbackRate: 1.0,
  addEventListener: vi.fn(),
  removeEventListener: vi.fn(),
})) as any;

describe('AudioPlayer', () => {
  beforeEach(() => {
    vi.clearAllMocks();
  });

  it('should render audio player', () => {
    render(<AudioPlayer audioUrl="https://example.com/audio.mp3" />);

    expect(screen.getByRole('button', { name: /再生/i })).toBeInTheDocument();
  });

  it('should play audio when play button is clicked', () => {
    render(<AudioPlayer audioUrl="https://example.com/audio.mp3" />);

    const playButton = screen.getByRole('button', { name: /再生/i });
    fireEvent.click(playButton);

    expect(mockPlay).toHaveBeenCalled();
  });

  it('should pause audio when pause button is clicked', () => {
    render(<AudioPlayer audioUrl="https://example.com/audio.mp3" />);

    const playButton = screen.getByRole('button', { name: /再生/i });
    fireEvent.click(playButton);

    const pauseButton = screen.getByRole('button', { name: /一時停止/i });
    fireEvent.click(pauseButton);

    expect(mockPause).toHaveBeenCalled();
  });

  it('should change speed', () => {
    render(<AudioPlayer audioUrl="https://example.com/audio.mp3" />);

    const speedButton = screen.getByRole('button', { name: /1.0x/i });
    fireEvent.click(speedButton);

    // 速度変更のメニューが表示される
    const speed15x = screen.getByRole('button', { name: /1.5x/i });
    fireEvent.click(speed15x);

    // 速度が変更される
    expect(screen.getByRole('button', { name: /1.5x/i })).toBeInTheDocument();
  });

  it('should show repeat button', () => {
    render(<AudioPlayer audioUrl="https://example.com/audio.mp3" />);

    expect(screen.getByRole('button', { name: /繰り返し/i })).toBeInTheDocument();
  });
});
