import { describe, it, expect } from 'vitest';
import { render, screen } from '@testing-library/react';
import Home from './page';

describe('Home', () => {
  it('renders heading', () => {
    render(<Home />);
    expect(screen.getByRole('heading', { name: /HaiLanGo/i })).toBeDefined();
  });

  it('renders description', () => {
    render(<Home />);
    expect(screen.getByText(/Coming Soon/i)).toBeDefined();
  });
});
