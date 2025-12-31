import type { Meta, StoryObj } from '@storybook/react';
import { Button } from './Button';

const meta = {
  title: 'UI/Button',
  component: Button,
  parameters: {
    layout: 'centered',
    docs: {
      description: {
        component:
          'Primary button component following HaiLanGo design system. Height: 48px (md), border-radius: 8px.',
      },
    },
  },
  tags: ['autodocs'],
  argTypes: {
    variant: {
      control: 'select',
      options: ['primary', 'secondary', 'outline', 'ghost', 'danger', 'success'],
      description: 'Button style variant',
    },
    size: {
      control: 'select',
      options: ['sm', 'md', 'lg'],
      description: 'Button size (sm: 36px, md: 48px, lg: 56px height)',
    },
    isLoading: {
      control: 'boolean',
      description: 'Shows loading spinner',
    },
    disabled: {
      control: 'boolean',
      description: 'Disables the button',
    },
    fullWidth: {
      control: 'boolean',
      description: 'Makes button full width',
    },
  },
} satisfies Meta<typeof Button>;

export default meta;
type Story = StoryObj<typeof meta>;

// Primary Button - Main action
export const Primary: Story = {
  args: {
    children: 'Primary Button',
    variant: 'primary',
    size: 'md',
  },
};

// Secondary Button - Secondary action
export const Secondary: Story = {
  args: {
    children: 'Secondary Button',
    variant: 'secondary',
    size: 'md',
  },
};

// Outline Button - Tertiary action
export const Outline: Story = {
  args: {
    children: 'Outline Button',
    variant: 'outline',
    size: 'md',
  },
};

// Ghost Button - Minimal style
export const Ghost: Story = {
  args: {
    children: 'Ghost Button',
    variant: 'ghost',
    size: 'md',
  },
};

// Danger Button - Destructive actions
export const Danger: Story = {
  args: {
    children: 'Delete',
    variant: 'danger',
    size: 'md',
  },
};

// Success Button - Positive actions
export const Success: Story = {
  args: {
    children: 'Complete',
    variant: 'success',
    size: 'md',
  },
};

// Size Variants
export const Small: Story = {
  args: {
    children: 'Small Button',
    variant: 'primary',
    size: 'sm',
  },
};

export const Medium: Story = {
  args: {
    children: 'Medium Button',
    variant: 'primary',
    size: 'md',
  },
};

export const Large: Story = {
  args: {
    children: 'Large Button',
    variant: 'primary',
    size: 'lg',
  },
};

// Loading State
export const Loading: Story = {
  args: {
    children: 'Loading',
    variant: 'primary',
    size: 'md',
    isLoading: true,
  },
};

// Disabled State
export const Disabled: Story = {
  args: {
    children: 'Disabled',
    variant: 'primary',
    size: 'md',
    disabled: true,
  },
};

// Full Width
export const FullWidth: Story = {
  args: {
    children: 'Full Width Button',
    variant: 'primary',
    size: 'md',
    fullWidth: true,
  },
  decorators: [
    (Story) => (
      <div style={{ width: '300px' }}>
        <Story />
      </div>
    ),
  ],
};

// With Left Icon
const PlayIcon = () => (
  <svg aria-hidden="true" fill="none" viewBox="0 0 24 24" stroke="currentColor">
    <path
      strokeLinecap="round"
      strokeLinejoin="round"
      strokeWidth={2}
      d="M14.752 11.168l-3.197-2.132A1 1 0 0010 9.87v4.263a1 1 0 001.555.832l3.197-2.132a1 1 0 000-1.664z"
    />
    <path
      strokeLinecap="round"
      strokeLinejoin="round"
      strokeWidth={2}
      d="M21 12a9 9 0 11-18 0 9 9 0 0118 0z"
    />
  </svg>
);

export const WithLeftIcon: Story = {
  args: {
    children: 'Play Audio',
    variant: 'primary',
    size: 'md',
    leftIcon: <PlayIcon />,
  },
};

// With Right Icon
const ArrowIcon = () => (
  <svg aria-hidden="true" fill="none" viewBox="0 0 24 24" stroke="currentColor">
    <path
      strokeLinecap="round"
      strokeLinejoin="round"
      strokeWidth={2}
      d="M13 7l5 5m0 0l-5 5m5-5H6"
    />
  </svg>
);

export const WithRightIcon: Story = {
  args: {
    children: 'Next Page',
    variant: 'primary',
    size: 'md',
    rightIcon: <ArrowIcon />,
  },
};

// All Variants Showcase
export const AllVariants: Story = {
  render: () => (
    <div className="flex flex-col gap-4">
      <div className="flex flex-wrap gap-3">
        <Button variant="primary">Primary</Button>
        <Button variant="secondary">Secondary</Button>
        <Button variant="outline">Outline</Button>
        <Button variant="ghost">Ghost</Button>
        <Button variant="danger">Danger</Button>
        <Button variant="success">Success</Button>
      </div>
      <div className="flex flex-wrap gap-3 items-center">
        <Button size="sm">Small</Button>
        <Button size="md">Medium</Button>
        <Button size="lg">Large</Button>
      </div>
      <div className="flex flex-wrap gap-3">
        <Button isLoading>Loading</Button>
        <Button disabled>Disabled</Button>
      </div>
    </div>
  ),
};

// HaiLanGo App Context Examples
export const LearningActions: Story = {
  render: () => (
    <div className="flex flex-col gap-4 p-6 bg-white rounded-xl shadow-soft">
      <h3 className="text-lg font-semibold text-text-primary">Learning Actions</h3>
      <div className="flex flex-col gap-3">
        <Button variant="primary" size="lg" fullWidth leftIcon={<PlayIcon />}>
          Start Learning
        </Button>
        <Button variant="secondary" size="md" fullWidth>
          Review Vocabulary
        </Button>
        <Button variant="outline" size="md" fullWidth>
          Skip This Page
        </Button>
      </div>
    </div>
  ),
  parameters: {
    backgrounds: { default: 'gray' },
  },
};

// Pronunciation Practice Buttons
const MicIcon = () => (
  <svg aria-hidden="true" fill="none" viewBox="0 0 24 24" stroke="currentColor">
    <path
      strokeLinecap="round"
      strokeLinejoin="round"
      strokeWidth={2}
      d="M19 11a7 7 0 01-7 7m0 0a7 7 0 01-7-7m7 7v4m0 0H8m4 0h4m-4-8a3 3 0 01-3-3V5a3 3 0 116 0v6a3 3 0 01-3 3z"
    />
  </svg>
);

export const PronunciationPractice: Story = {
  render: () => (
    <div className="flex flex-col gap-4 p-6 bg-white rounded-xl shadow-soft max-w-sm">
      <h3 className="text-lg font-semibold text-text-primary">Pronunciation Practice</h3>
      <p className="text-text-secondary text-sm">Press the button and speak clearly</p>
      <Button variant="danger" size="lg" fullWidth leftIcon={<MicIcon />}>
        Start Recording
      </Button>
      <div className="flex gap-2">
        <Button variant="ghost" size="sm">
          Skip
        </Button>
        <Button variant="outline" size="sm" className="flex-1">
          Play Example
        </Button>
      </div>
    </div>
  ),
  parameters: {
    backgrounds: { default: 'gray' },
  },
};
