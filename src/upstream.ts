export const getUpstreamUrl = (path: string) => {
  if (path.includes("/sqs")) {
    return `${process.env.SENTINEL_EXTERNAL_URL}/sqs`;
  }

  return process.env.DESTINATION_URL;
};
