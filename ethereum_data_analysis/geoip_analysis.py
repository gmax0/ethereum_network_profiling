import geoip2.database
import argparse
import csv

parser = argparse.ArgumentParser(description='Process some integers.')
parser.add_argument('--city_db_location', help='full path to the city .mmdb')
parser.add_argument('--asn_db_location', help='full path to the asn .mmdb')
parser.add_argument('--peerfile_location', help="full path to list of peer IPs")

args = parser.parse_args()

with open(args.peerfile_location) as f:
    peers = f.read().splitlines() 

cols = ['ip', 'asn', 'asorg', 'latitude', 'longitude', 'city', 'country']

with geoip2.database.Reader(args.city_db_location) as city_reader, geoip2.database.Reader(args.asn_db_location) as asn_reader, open('enriched_peers.csv', 'w') as csv_file:
    writer = csv.writer(csv_file)
    writer.writerow(cols)

    for peer in peers:
        peer_data = []

        # Add IP
        peer_data.append(peer)

        # Grab ASN
        response = asn_reader.asn(peer)
        peer_data.append(response.autonomous_system_number)
        peer_data.append(response.autonomous_system_organization)

        # Add lattitude + longitude
        response = city_reader.city(peer)
        peer_data.append(response.location.latitude)
        peer_data.append(response.location.longitude)
        peer_data.append(response.city.name)
        peer_data.append(response.country.name)

        writer.writerow(peer_data)