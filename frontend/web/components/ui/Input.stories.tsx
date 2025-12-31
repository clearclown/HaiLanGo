import type { Meta, StoryObj } from '@storybook/react';
import { Button } from './Button';
import { Input, Textarea } from './Input';

const meta = {
  title: 'UI/Input',
  component: Input,
  parameters: {
    layout: 'centered',
    docs: {
      description: {
        component:
          'Input component following HaiLanGo design system. Height: 48px (md), border-radius: 8px.',
      },
    },
  },
  tags: ['autodocs'],
  argTypes: {
    variant: {
      control: 'select',
      options: ['default', 'filled'],
      description: 'Input style variant',
    },
    inputSize: {
      control: 'select',
      options: ['sm', 'md', 'lg'],
      description: 'Input size (sm: 36px, md: 48px, lg: 56px height)',
    },
    disabled: {
      control: 'boolean',
      description: 'Disables the input',
    },
  },
  decorators: [
    (Story) => (
      <div style={{ width: '320px' }}>
        <Story />
      </div>
    ),
  ],
} satisfies Meta<typeof Input>;

export default meta;
type Story = StoryObj<typeof meta>;

// Default Input
export const Default: Story = {
  args: {
    placeholder: 'Enter your text...',
  },
};

// With Label
export const WithLabel: Story = {
  args: {
    label: 'Email',
    placeholder: 'your@email.com',
    type: 'email',
  },
};

// With Hint
export const WithHint: Story = {
  args: {
    label: 'Username',
    placeholder: 'Choose a username',
    hint: 'This will be your public display name.',
  },
};

// With Error
export const WithError: Story = {
  args: {
    label: 'Email',
    placeholder: 'your@email.com',
    value: 'invalid-email',
    error: 'Please enter a valid email address.',
  },
};

// Password Input
export const Password: Story = {
  args: {
    label: 'Password',
    placeholder: 'Enter your password',
    type: 'password',
  },
};

// With Left Icon
const SearchIcon = () => (
  <svg aria-hidden="true" fill="none" viewBox="0 0 24 24" stroke="currentColor">
    <path
      strokeLinecap="round"
      strokeLinejoin="round"
      strokeWidth={2}
      d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z"
    />
  </svg>
);

export const WithLeftIcon: Story = {
  args: {
    placeholder: 'Search...',
    leftIcon: <SearchIcon />,
  },
};

// Filled Variant
export const Filled: Story = {
  args: {
    label: 'Search Books',
    placeholder: 'Enter book title or author',
    variant: 'filled',
    leftIcon: <SearchIcon />,
  },
};

// All Sizes
export const AllSizes: Story = {
  render: () => (
    <div className="space-y-4">
      <Input inputSize="sm" placeholder="Small input" />
      <Input inputSize="md" placeholder="Medium input" />
      <Input inputSize="lg" placeholder="Large input" />
    </div>
  ),
};

// Disabled State
export const Disabled: Story = {
  args: {
    label: 'Email',
    placeholder: 'Cannot edit this',
    disabled: true,
    value: 'disabled@example.com',
  },
};

// Login Form Example
export const LoginForm: Story = {
  render: () => (
    <div className="p-6 bg-white rounded-xl shadow-soft space-y-4">
      <h2 className="text-xl font-semibold text-text-primary mb-6">Login</h2>
      <Input
        label="Email"
        placeholder="your@email.com"
        type="email"
        leftIcon={
          <svg aria-hidden="true" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path
              strokeLinecap="round"
              strokeLinejoin="round"
              strokeWidth={2}
              d="M16 12a4 4 0 10-8 0 4 4 0 008 0zm0 0v1.5a2.5 2.5 0 005 0V12a9 9 0 10-9 9m4.5-1.206a8.959 8.959 0 01-4.5 1.207"
            />
          </svg>
        }
      />
      <Input label="Password" placeholder="Enter your password" type="password" />
      <div className="flex items-center justify-between">
        <label className="flex items-center gap-2 cursor-pointer">
          <input
            type="checkbox"
            className="w-4 h-4 rounded border-border text-primary focus:ring-primary"
          />
          <span className="text-sm text-text-secondary">Remember me</span>
        </label>
        <button type="button" className="text-sm text-primary hover:underline">
          Forgot password?
        </button>
      </div>
      <Button variant="primary" fullWidth>
        Sign In
      </Button>
    </div>
  ),
  parameters: {
    backgrounds: { default: 'gray' },
  },
};

// Search Input
export const SearchInput: Story = {
  render: () => (
    <div className="p-4 bg-white rounded-xl shadow-soft">
      <Input
        variant="filled"
        placeholder="Search books, phrases, vocabulary..."
        leftIcon={<SearchIcon />}
        inputSize="lg"
      />
    </div>
  ),
  parameters: {
    backgrounds: { default: 'gray' },
  },
};

// Book Title Input (HaiLanGo specific)
export const BookTitleInput: Story = {
  render: () => (
    <div className="p-6 bg-white rounded-xl shadow-soft space-y-4">
      <h3 className="font-semibold text-text-primary">Add New Book</h3>
      <Input
        label="Book Title"
        placeholder="e.g., Russian Language 101"
        hint="This will be displayed in your bookshelf"
      />
      <div className="grid grid-cols-2 gap-4">
        <div>
          <label
            htmlFor="target-language"
            className="block text-sm font-medium text-text-primary mb-2"
          >
            Target Language
          </label>
          <select
            id="target-language"
            className="w-full h-12 px-4 rounded-lg border border-border bg-white text-text-primary focus:border-primary focus:ring-2 focus:ring-primary/20 focus:outline-none"
          >
            <option>Russian</option>
            <option>Arabic</option>
            <option>Persian</option>
            <option>Hebrew</option>
          </select>
        </div>
        <div>
          <label
            htmlFor="native-language"
            className="block text-sm font-medium text-text-primary mb-2"
          >
            Native Language
          </label>
          <select
            id="native-language"
            className="w-full h-12 px-4 rounded-lg border border-border bg-white text-text-primary focus:border-primary focus:ring-2 focus:ring-primary/20 focus:outline-none"
          >
            <option>Japanese</option>
            <option>English</option>
            <option>Chinese</option>
          </select>
        </div>
      </div>
    </div>
  ),
  parameters: {
    backgrounds: { default: 'gray' },
  },
};

// Textarea Stories
export const TextareaDefault: StoryObj<typeof Textarea> = {
  render: () => <Textarea label="Notes" placeholder="Add your learning notes here..." rows={4} />,
};

export const TextareaWithError: StoryObj<typeof Textarea> = {
  render: () => (
    <Textarea
      label="Feedback"
      placeholder="Share your thoughts..."
      rows={4}
      error="Please enter at least 10 characters."
      value="Too short"
    />
  ),
};

export const TextareaFilled: StoryObj<typeof Textarea> = {
  render: () => (
    <Textarea
      label="Translation Notes"
      placeholder="Add context or translation notes for this phrase..."
      rows={3}
      variant="filled"
    />
  ),
};

// OCR Correction Input (HaiLanGo specific)
export const OCRCorrectionForm: Story = {
  render: () => (
    <div className="p-6 bg-white rounded-xl shadow-soft space-y-4 w-96">
      <h3 className="font-semibold text-text-primary">Correct OCR Text</h3>
      <p className="text-sm text-text-secondary">
        If the OCR result is incorrect, you can manually correct it below.
      </p>
      <div className="p-4 bg-background-secondary rounded-lg">
        <p className="text-sm text-text-secondary mb-1">Detected text:</p>
        <p className="text-text-primary font-medium">3дравствуйте!</p>
      </div>
      <Input
        label="Corrected Text"
        placeholder="Enter the correct text"
        defaultValue="Здравствуйте!"
      />
      <Textarea
        label="Translation"
        placeholder="Add translation (optional)"
        rows={2}
        defaultValue="Hello! (formal)"
      />
      <div className="flex gap-3">
        <Button variant="outline" className="flex-1">
          Cancel
        </Button>
        <Button variant="primary" className="flex-1">
          Save Correction
        </Button>
      </div>
    </div>
  ),
  parameters: {
    backgrounds: { default: 'gray' },
  },
};
