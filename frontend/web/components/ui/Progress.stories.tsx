import type { Meta, StoryObj } from '@storybook/react';
import { Progress, ProgressCircle } from './Progress';

const meta = {
  title: 'UI/Progress',
  component: Progress,
  parameters: {
    layout: 'centered',
    docs: {
      description: {
        component:
          'Progress bar component for showing completion status. Used for learning progress, upload status, etc.',
      },
    },
  },
  tags: ['autodocs'],
  argTypes: {
    value: {
      control: { type: 'range', min: 0, max: 100 },
      description: 'Current progress value',
    },
    variant: {
      control: 'select',
      options: ['primary', 'secondary', 'success', 'warning', 'danger', 'gradient'],
      description: 'Progress bar color variant',
    },
    size: {
      control: 'select',
      options: ['sm', 'md', 'lg'],
      description: 'Progress bar height',
    },
    showLabel: {
      control: 'boolean',
      description: 'Show percentage label',
    },
    animated: {
      control: 'boolean',
      description: 'Animate the progress bar',
    },
  },
  decorators: [
    (Story) => (
      <div style={{ width: '300px' }}>
        <Story />
      </div>
    ),
  ],
} satisfies Meta<typeof Progress>;

export default meta;
type Story = StoryObj<typeof meta>;

// Default Progress
export const Default: Story = {
  args: {
    value: 45,
    variant: 'primary',
    size: 'md',
  },
};

// With Label
export const WithLabel: Story = {
  args: {
    value: 65,
    variant: 'primary',
    size: 'md',
    showLabel: true,
  },
};

// Custom Label
export const CustomLabel: Story = {
  args: {
    value: 3,
    max: 5,
    variant: 'primary',
    size: 'md',
    showLabel: true,
    label: 'Pages completed',
  },
};

// Gradient Variant
export const Gradient: Story = {
  args: {
    value: 75,
    variant: 'gradient',
    size: 'md',
    showLabel: true,
    label: 'Learning progress',
  },
};

// All Variants
export const AllVariants: Story = {
  render: () => (
    <div className="space-y-4 w-72">
      <div>
        <p className="text-sm text-text-secondary mb-2">Primary</p>
        <Progress value={60} variant="primary" />
      </div>
      <div>
        <p className="text-sm text-text-secondary mb-2">Secondary</p>
        <Progress value={60} variant="secondary" />
      </div>
      <div>
        <p className="text-sm text-text-secondary mb-2">Success</p>
        <Progress value={60} variant="success" />
      </div>
      <div>
        <p className="text-sm text-text-secondary mb-2">Warning</p>
        <Progress value={60} variant="warning" />
      </div>
      <div>
        <p className="text-sm text-text-secondary mb-2">Danger</p>
        <Progress value={60} variant="danger" />
      </div>
      <div>
        <p className="text-sm text-text-secondary mb-2">Gradient</p>
        <Progress value={60} variant="gradient" />
      </div>
    </div>
  ),
};

// All Sizes
export const AllSizes: Story = {
  render: () => (
    <div className="space-y-4 w-72">
      <div>
        <p className="text-sm text-text-secondary mb-2">Small</p>
        <Progress value={60} size="sm" variant="primary" />
      </div>
      <div>
        <p className="text-sm text-text-secondary mb-2">Medium</p>
        <Progress value={60} size="md" variant="primary" />
      </div>
      <div>
        <p className="text-sm text-text-secondary mb-2">Large</p>
        <Progress value={60} size="lg" variant="primary" />
      </div>
    </div>
  ),
};

// Animated Progress
export const Animated: Story = {
  args: {
    value: 50,
    variant: 'primary',
    size: 'md',
    animated: true,
    showLabel: true,
    label: 'Processing...',
  },
};

// Learning Progress Example
export const LearningProgress: Story = {
  render: () => (
    <div className="space-y-6 w-72">
      <div>
        <h3 className="text-sm font-medium text-text-primary mb-3">Book Progress</h3>
        <div className="space-y-3">
          <Progress value={45} variant="gradient" showLabel label="Russian 101" />
          <Progress value={78} variant="gradient" showLabel label="Arabic Basics" />
          <Progress value={12} variant="gradient" showLabel label="Persian Grammar" />
        </div>
      </div>
    </div>
  ),
};

// Circle Progress Meta
const circleMeta = {
  title: 'UI/ProgressCircle',
  component: ProgressCircle,
  parameters: {
    layout: 'centered',
    docs: {
      description: {
        component: 'Circular progress indicator for scores and achievements.',
      },
    },
  },
  tags: ['autodocs'],
  argTypes: {
    value: {
      control: { type: 'range', min: 0, max: 100 },
    },
    variant: {
      control: 'select',
      options: ['primary', 'secondary', 'success', 'warning', 'danger'],
    },
    size: {
      control: 'select',
      options: ['sm', 'md', 'lg'],
    },
  },
} satisfies Meta<typeof ProgressCircle>;

// Circle Progress Stories
export const CircleDefault: StoryObj<typeof ProgressCircle> = {
  render: () => <ProgressCircle value={75} variant="primary" size="md" />,
};

export const CircleAllSizes: StoryObj<typeof ProgressCircle> = {
  render: () => (
    <div className="flex items-center gap-6">
      <ProgressCircle value={75} size="sm" />
      <ProgressCircle value={75} size="md" />
      <ProgressCircle value={75} size="lg" />
    </div>
  ),
};

export const CircleAllVariants: StoryObj<typeof ProgressCircle> = {
  render: () => (
    <div className="flex items-center gap-6">
      <ProgressCircle value={85} variant="primary" />
      <ProgressCircle value={70} variant="secondary" />
      <ProgressCircle value={92} variant="success" />
      <ProgressCircle value={45} variant="warning" />
      <ProgressCircle value={30} variant="danger" />
    </div>
  ),
};

// Pronunciation Score Example
export const PronunciationScore: StoryObj<typeof ProgressCircle> = {
  render: () => (
    <div className="flex flex-col items-center gap-4 p-6 bg-white rounded-xl shadow-soft">
      <ProgressCircle value={85} variant="success" size="lg" strokeWidth={6} />
      <div className="text-center">
        <h3 className="font-semibold text-text-primary">Great Job!</h3>
        <p className="text-sm text-text-secondary">Your pronunciation score</p>
      </div>
    </div>
  ),
  parameters: {
    backgrounds: { default: 'gray' },
  },
};

// Stats Dashboard
export const StatsDashboard: StoryObj<typeof ProgressCircle> = {
  render: () => (
    <div className="grid grid-cols-3 gap-6 p-6 bg-white rounded-xl shadow-soft">
      <div className="flex flex-col items-center gap-2">
        <ProgressCircle value={85} variant="primary" size="md" />
        <span className="text-xs text-text-secondary">Accuracy</span>
      </div>
      <div className="flex flex-col items-center gap-2">
        <ProgressCircle value={72} variant="secondary" size="md" />
        <span className="text-xs text-text-secondary">Fluency</span>
      </div>
      <div className="flex flex-col items-center gap-2">
        <ProgressCircle value={90} variant="success" size="md" />
        <span className="text-xs text-text-secondary">Completion</span>
      </div>
    </div>
  ),
  parameters: {
    backgrounds: { default: 'gray' },
  },
};

// Combined Progress Example
export const CombinedProgress: StoryObj<typeof Progress> = {
  render: () => (
    <div className="p-6 bg-white rounded-xl shadow-soft w-80">
      <h3 className="font-semibold text-text-primary mb-4">Today&apos;s Progress</h3>
      <div className="flex items-center gap-6 mb-6">
        <ProgressCircle value={65} variant="primary" size="lg" />
        <div>
          <p className="text-2xl font-bold text-text-primary">65%</p>
          <p className="text-sm text-text-secondary">Overall completion</p>
        </div>
      </div>
      <div className="space-y-3">
        <Progress value={80} variant="primary" showLabel label="Vocabulary" size="sm" />
        <Progress value={65} variant="secondary" showLabel label="Reading" size="sm" />
        <Progress value={50} variant="warning" showLabel label="Speaking" size="sm" />
      </div>
    </div>
  ),
  parameters: {
    backgrounds: { default: 'gray' },
  },
};
