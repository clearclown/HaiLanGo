import { render, screen } from '@testing-library/react';
import { describe, expect, it } from 'vitest';
import { WelcomeCard } from './WelcomeCard';

describe('WelcomeCard', () => {
  it('should render user name', () => {
    render(<WelcomeCard userName="太郎" />);

    expect(screen.getByText(/こんにちは、太郎さん/)).toBeDefined();
  });

  it('should display motivation message', () => {
    render(<WelcomeCard userName="太郎" />);

    expect(screen.getByText(/今日も頑張りましょう！/)).toBeDefined();
  });

  it('should display greeting icon', () => {
    const { container } = render(<WelcomeCard userName="太郎" />);

    // Check for wave emoji or icon
    expect(container.textContent).toContain('👋');
  });
});
