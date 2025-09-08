# Description: GraphQL queries for the AIB API


def get_all_cameras() -> str:
    """
    Get all available cameras query
    :return: Cameras formated query
    """
    query = """
    query GetCameras {
        cameras {
            id
            name
            communicationStatus {
                started
                failing
            }
            videoStreams {
                id
                name
                videoCodec
                streamAvailability{
                    rtsp
                }
            }
        }
    }"""
    
    return query.strip()


def get_about():
    """
    Get About query
    :return: About formated query
    """
    query = """
    query about {
      about {
        videoManagementSystems{
            id
            url
            idp
            version
            vendorID
            slc
        }
      }
    }"""

    return query.strip()

