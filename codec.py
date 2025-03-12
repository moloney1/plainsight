import os
import shutil

from PIL import Image

def message_to_binary(message: str) -> str:
    msg = [f"{bin(ord(character))[2:]:0>8}" for character in message]
    return "".join(msg)


def message_from_binary(message: str) -> str:
    print(f"{message=}")
    decoded_message = []
    bytes_read = 0
    start = 0
    while bytes_read < int(len(message) / 8):
        decoded_message.append(chr(int(message[start:start+8], 2)))
        print(f"{decoded_message=}")
        start += 8
        bytes_read += 1

    return "".join(decoded_message)


def modify_pixel(im: Image, coords: tuple, new_pixel: tuple) -> None:
    """
    e.g.  modify_pixel(im, (0, 0), (255, 255, 255)) 
    """
    im.putpixel(coords, new_pixel)


def write_message_to_image(image: Image, message: str, output_filename: str = "out.png") -> None:
    """
    1. copy image data to memory
    2. iterate and modify the bytes to encode message
    """
    # add an alpha layer if not there already
    #if image.mode == "RGB":
    #    image.putalpha(1)
    encoded_msg = message_to_binary(message)
    encoded_msg_pointer = 0

    new_image = []
    for p in image.getdata():
        if encoded_msg_pointer >= len(encoded_msg):
            new_image.append(p)
        else:
            new_pixel = []
            for frame in p:
                if encoded_msg_pointer < len(encoded_msg):
                    frame_as_bin = list(format(frame, "08b"))
                    frame_as_bin[-1] = encoded_msg[encoded_msg_pointer]
                    frame_as_str = "".join(str(bit) for bit in frame_as_bin)
                    new_pixel.append(int(frame_as_str, 2))
                    encoded_msg_pointer += 1
                else:
                    new_pixel.append(frame)
            new_image.append(tuple(new_pixel))

    print(encoded_msg)
    print(new_image[:10])
    ni = Image.new("RGB", image.size)
    ni.putdata(new_image)
    ni.save(output_filename)


def read_message_from_image(image: Image, bytes_to_read: int) -> str:
    count = 0
    limit = 8 * bytes_to_read
    message = []
    for p in image.getdata():
        for f in p:
            if count < limit:
                frame_as_bin = list(format(f, "08b"))
                message.append(frame_as_bin[-1])
                count += 1
            else:
                return message_from_binary("".join(message))


def main():
    if not os.path.exists(out_path := os.path.join(os.getcwd(), "output")):
        os.mkdir(out_path)

    shutil.copy(
        os.path.join(os.getcwd(), "assets", "test_img.png"),
        os.path.join(os.getcwd(), "output")
        )

    im = Image.open(os.path.join(os.getcwd(), "output", "test_img.png"))


    # test 1
#    new_image_data = list(im.getdata())
#    new_image_data[0] = (255, 255, 255)
#    new_image_data[1] = (255, 255, 255)
#    new_image_data[3] = (255, 255, 255)
#    new_image_data[7] = (255, 255, 255)
#    new_image_data[252] = (255, 255, 255)
#    new_image_data[377] = (255, 255, 255)
#    new_image_data[1234] = (255, 255, 255)
#    new_image_data[5678] = (255, 255, 255)
#    new_image = Image.new("RGB", im.size)
#    new_image.putdata(new_image_data)
#    new_image.show()

    # test 2
#    new_image_data = []
#    for i, pix in enumerate(im.getdata()):
#        r, g, b, a = pix
#        print(format(r, '08b'))
#    new_image = Image.new("RGBA", im.size)
#    new_image.putdata(new_image_data)
#    new_image.show()

    write_message_to_image(im, "hi")
    print("hi in binary:")
    print(message_to_binary("hi"))

    out_im = Image.open(os.path.join(os.getcwd(), "out.png"))
    m = read_message_from_image(out_im, 2)
    print(m)


if __name__ == "__main__":
    main()
