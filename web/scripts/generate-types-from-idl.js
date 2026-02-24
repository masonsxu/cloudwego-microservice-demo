#!/usr/bin/env node

/**
 * Thrift IDL to TypeScript Type Generator
 *
 * 从 Thrift IDL 文件生成 TypeScript 类型定义
 * 确保前后端类型定义一致性
 */

const fs = require('fs')
const path = require('path')

const IDL_DIR = path.resolve(__dirname, '../../idl')
const OUTPUT_DIR = path.resolve(__dirname, '../src/types/idl')

// 类型映射
const typeMapping = {
  'string': 'string',
  'bool': 'boolean',
  'i8': 'number',
  'i16': 'number',
  'i32': 'number',
  'i64': 'number',
  'double': 'number',
  'binary': 'string',
  'list': 'Array',
  'set': 'Set',
  'map': 'Record'
}

/**
 * 解析 Thrift struct 定义
 */
function parseStruct(content) {
  const structRegex = /struct\s+(\w+)\s*{([^}]+)}/gs
  const structs = {}
  let match

  while ((match = structRegex.exec(content)) !== null) {
    const structName = match[1]
    const structBody = match[2]
    const fields = []

    // 解析字段
    const fieldRegex = /(\d+):\s*(?:optional|required)\s+(\w+)\s+(\w+)(?:\s*\(go\.tag\s*=\s*"([^"]+)"\))?/g
    let fieldMatch

    while ((fieldMatch = fieldRegex.exec(structBody)) !== null) {
      const fieldId = fieldMatch[1]
      const required = fieldMatch[2]
      const fieldType = fieldMatch[3]
      const fieldName = fieldMatch[4]
      const jsonTag = fieldMatch[5] || fieldName

      fields.push({
        id: fieldId,
        required: required === 'required',
        type: fieldType,
        name: fieldName,
        jsonTag: jsonTag
      })
    }

    structs[structName] = {
      name: structName,
      fields
    }
  }

  return structs
}

/**
 * 转换 Thrift 类型为 TypeScript 类型
 */
function convertType(thriftType) {
  // 处理泛型类型
  const genericMatch = thriftType.match(/^(list|set|map)<(.+)>$/)
  if (genericMatch) {
    const genericType = genericMatch[1]
    const innerType = genericMatch[2]

    if (genericType === 'map') {
      const types = innerType.split(',').map(t => convertType(t.trim()))
      return `Record<${types[0]}, ${types[1]}>`
    } else {
      const inner = convertType(innerType)
      return `${typeMapping[genericType]}<${inner}>`
    }
  }

  // 处理普通类型
  return typeMapping[thriftType] || thriftType
}

/**
 * 生成 TypeScript 接口
 */
function generateTypeScriptInterface(struct) {
  const lines = []
  
  lines.push(`export interface ${struct.name} {`)
  
  struct.fields.forEach(field => {
    const optional = field.required ? '' : '?'
    const tsType = convertType(field.type)
    const comment = field.required ? '' : ' // 可选'
    lines.push(`  /** ${comment} */`)
    lines.push(`  ${field.name}${optional}: ${tsType}`)
  })
  
  lines.push('}')
  
  return lines.join('\n')
}

/**
 * 处理单个文件
 */
function processFile(inputFile, outputFile) {
  console.log(`处理文件: ${inputFile}`)
  
  const content = fs.readFileSync(inputFile, 'utf-8')
  const structs = parseStruct(content)
  
  const output = []
  output.push('/**')
  output.push(` * 从 ${path.basename(inputFile)} 自动生成的 TypeScript 类型定义`)
  output.push(' * 请勿手动修改')
  output.push(' */')
  output.push('')
  
  Object.values(structs).forEach(struct => {
    output.push(generateTypeScriptInterface(struct))
    output.push('')
  })
  
  fs.writeFileSync(outputFile, output.join('\n'))
  console.log(`✓ 生成: ${outputFile}`)
}

/**
 * 主函数
 */
function main() {
  console.log('开始从 Thrift IDL 生成 TypeScript 类型...\n')
  
  // 确保输出目录存在
  if (!fs.existsSync(OUTPUT_DIR)) {
    fs.mkdirSync(OUTPUT_DIR, { recursive: true })
  }
  
  // 需要处理的文件列表
  const filesToProcess = [
    {
      input: path.join(IDL_DIR, 'http/identity/identity_model.thrift'),
      output: path.join(OUTPUT_DIR, 'identity_model.ts')
    },
    {
      input: path.join(IDL_DIR, 'http/base/base.thrift'),
      output: path.join(OUTPUT_DIR, 'base.ts')
    },
    {
      input: path.join(IDL_DIR, 'http/permission/permission_model.thrift'),
      output: path.join(OUTPUT_DIR, 'permission_model.ts')
    }
  ]
  
  let successCount = 0
  let errorCount = 0
  
  filesToProcess.forEach(({ input, output }) => {
    try {
      if (fs.existsSync(input)) {
        processFile(input, output)
        successCount++
      } else {
        console.warn(`⚠ 文件不存在: ${input}`)
        errorCount++
      }
    } catch (error) {
      console.error(`✗ 处理失败: ${input}`)
      console.error(error)
      errorCount++
    }
  })
  
  console.log(`\n完成! 成功: ${successCount}, 失败: ${errorCount}`)
  
  if (successCount > 0) {
    console.log(`\n生成的文件保存在: ${OUTPUT_DIR}`)
  }
}

// 运行
main()
