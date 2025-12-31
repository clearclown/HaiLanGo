import { fireEvent, render, screen } from '@testing-library/react';
import { beforeEach, describe, expect, it, vi } from 'vitest';
import { AudioPlayer } from './AudioPlayer';

// useAudioPlayer hookをモック
const mockTogglePlayPause = vi.fn();
const mockSetSpeed = vi.fn();

vi.mock('@/hooks/useAudioPlayer', () => ({
  useAudioPlayer: () => ({
    playing: false,
    currentTime: 30,
    duration: 100,
    speed: 1.0,
    setSpeed: mockSetSpeed,
    togglePlayPause: mockTogglePlayPause,
  }),
}));

describe('AudioPlayer', () => {
  beforeEach(() => {
    vi.clearAllMocks();
  });

  it('should render audio player', () => {
    render(<AudioPlayer audioUrl="https://example.com/audio.mp3" />);

    expect(screen.getByRole('button', { name: /再生/i })).toBeInTheDocument();
  });

  it('should call togglePlayPause when play button is clicked', () => {
    render(<AudioPlayer audioUrl="https://example.com/audio.mp3" />);

    const playButton = screen.getByRole('button', { name: /再生/i });
    fireEvent.click(playButton);

    expect(mockTogglePlayPause).toHaveBeenCalled();
  });

  it('should show speed menu when speed button is clicked', () => {
    render(<AudioPlayer audioUrl="https://example.com/audio.mp3" />);

    // Click the speed button (labeled "1.0x")
    const speedButton = screen.getByRole('button', { name: '1x' });
    fireEvent.click(speedButton);

    // Speed menu should appear with various options
    expect(screen.getByRole('button', { name: '0.5x' })).toBeInTheDocument();
    expect(screen.getByRole('button', { name: '1.5x' })).toBeInTheDocument();
    expect(screen.getByRole('button', { name: '2x' })).toBeInTheDocument();
  });

  it('should call setSpeed when speed option is selected', () => {
    render(<AudioPlayer audioUrl="https://example.com/audio.mp3" />);

    // Open speed menu
    const speedButton = screen.getByRole('button', { name: '1x' });
    fireEvent.click(speedButton);

    // Click 1.5x option
    const speed15x = screen.getByRole('button', { name: '1.5x' });
    fireEvent.click(speed15x);

    expect(mockSetSpeed).toHaveBeenCalledWith(1.5);
  });

  it('should show repeat button', () => {
    render(<AudioPlayer audioUrl="https://example.com/audio.mp3" />);

    expect(screen.getByRole('button', { name: /繰り返し/i })).toBeInTheDocument();
  });

  it('should show current time and duration', () => {
    render(<AudioPlayer audioUrl="https://example.com/audio.mp3" />);

    expect(screen.getByText('0:30 / 1:40')).toBeInTheDocument();
  });
});
