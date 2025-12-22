// Package primitives provides low-level layout primitives for ForgeUI
//
// Primitives are the building blocks for more complex components.
// They provide type-safe wrappers around common CSS layout patterns.
//
// All primitives only import from the root forgeui package to avoid
// circular dependencies. They should not import from other component packages.
package primitives

// Re-export commonly used primitive functions for convenience
// This allows users to import just "primitives" and access all primitives

// Layout primitives:
// - Box: Polymorphic container element
// - Flex: Flexbox container
// - Grid: CSS Grid container

// Stack helpers:
// - VStack: Vertical stack (flex column)
// - HStack: Horizontal stack (flex row)

// Utility primitives:
// - Center: Centers content both horizontally and vertically
// - Container: Responsive container with max-width
// - Spacer: Flexible spacer that fills available space
// - Text: Typography primitive
