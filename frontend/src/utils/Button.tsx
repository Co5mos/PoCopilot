import styled from '@emotion/styled'

const Button = styled('button')`
  padding: 20px;
  background-color: ${(props: { theme: { someLibProperty: any } }) => props.theme.someLibProperty};
  border-radius: 3px;
`

export default Button