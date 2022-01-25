import Icon from "@chakra-ui/icon"
import { Link, HStack, Text, Box, Flex } from "@chakra-ui/layout"
import { VscGithub, VscTwitter } from "react-icons/vsc"
import { Card } from "../../atoms/Card"
import { LogoIcon } from "../../atoms/LogoIcon"

export function About() {
  return (
    <>
      <Card title="Creators">
        <Flex direction={{ base: "column", md: "column", lg: "row" }} gridGap={{ base: 4, md: 4, lg: 8 }} >
          <Box>
            <Text variant="brandSecondary">Team: </Text>
            <Link href="https://twitter.com/ubiqsmart">
              <HStack>
                <Icon as={VscTwitter} title="Follow Ubiq on Twitter" />
                <Text>@ubiqsmart</Text>
              </HStack>
            </Link>
          </Box>

          <Box>
            <Text variant="brandSecondary">Open Sourced: </Text>
            <Link href="https://github.com/ubiq/ubiq-burn-stats" target="_blank">
              <HStack>
                <Icon as={VscGithub} title="Contribute to Ubiq Burn Stats on GitHub" />
                <Text>ubiq-burn-stats</Text>
              </HStack>
            </Link>
          </Box>

          <Box>
            <Text variant="brandSecondary">Tech Stack: </Text>
            <Text>Golang, TypeScript, Chakra</Text>
          </Box>
        </Flex>
      </Card>

      <Card title="About">
        <Text>
          Ethereum <Link href="https://eips.ethereum.org/EIPS/eip-1559" target="_blank">EIP-1559</Link> is an improvement proposal that
          makes changes to the transaction fee system. EIP-1559 got rid of the first-price auction which was a major source of transaction
          fees and replaced it with the base fee model where the fee is changed dynamically based on network activity. The base fee is
          burned <LogoIcon /> and not given to the miners.
        </Text>
        <Text>This website showcases how much has been burned among other neat stats. For more information, Vitalik Buterin {" "}
          <Link href="https://notes.ethereum.org/@vbuterin/eip-1559-faq" target="_blank">wrote a great blog post</Link> explaining
          what it is in detail.
        </Text>
      </Card>
    </>
  )
}
