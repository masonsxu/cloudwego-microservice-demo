/* generated using openapi-typescript-codegen -- do not edit */
/* istanbul ignore file */
/* tslint:disable */
/* eslint-disable */

// Permission levels (使用对象替代 enum 以支持 erasableSyntaxOnly)
export const PermissionLevel = {
  NONE: 0,
  READ: 1,
  WRITE: 2,
  FULL: 3,
} as const

export type PermissionLevel = typeof PermissionLevel[keyof typeof PermissionLevel]
